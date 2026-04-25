// boot/types/boot_modes.go

package boot_types

import (
	"errors"
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot"
	boot_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/phases"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	security_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	security_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/verification"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/keys"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

// DecideBootPath determines whether to run fast or cold boot
func (bm *boot_phase.BootManager) DecideBootPath() (*schema_system.BootSequence, error) {

	if bm.Identity == nil || bm.Identity.MachineID == "" {
		return nil, errors.New("invalid machine identity")
	}

	lastEnvKey := keys.LastKnownEnvKey(bm.Identity.MachineID)

	env, err := bm.Vault.LoadConfig(lastEnvKey)
	if err != nil {
		return bm.runColdBoot(bm.Identity)
	}

	if env.SchemaVersion < schema_system.CurrentVersion {
		env = schema_system.Migrate(env)
		if err := bm.Vault.SaveConfig(lastEnvKey, env); err != nil {
			return nil, err
		}
	}

	if _, err := bm.Vault.LoadGoldenHash(bm.Identity.MachineID); err != nil {
		return bm.runColdBoot(bm.Identity)
	}

	return bm.runFastBoot(env)
}

// ------------------------------------------------------------
// Cold Boot: full hardware discovery and provisioning
// ------------------------------------------------------------
func (bm *boot_phase.BootManager) runColdBoot(identity *schema_system.MachineIdentity) (*schema_system.BootSequence, error) {
	// Use the discovered platform + identity info
	env := &schema_system.EnvConfig{
		Platform: schema_system.PlatformResolution{
			Final:      identity.PlatformType, // << propagate platform
			Locked:     false,
			Source:     "discovery",
			ResolvedAt: time.Now(),
		},
		Identity: *identity,
		Hardware: identity.Hardware,
	}

	logging.Info(
		"[func (bm *BootManager) runColdBoot()] Platform: %s | OS: %s | Arch: %s | EntityType: %v",
		env.Platform.Final,
		env.Identity.OS,
		env.Identity.Arch,
		env.Identity.EntityType,
	)

	fullProfile, err := probe.ActiveDiscovery(env)
	if err != nil {
		return nil, fmt.Errorf("hardware discovery failed: %w", err)
	}

	// Update identity with discovered hardware
	bm.Identity.BindHardware(fullProfile)

	// 2. Provision golden baseline
	goldenHash, err := security_verification.ProvisionGolden(bm.Vault, bm.Identity.MachineID)
	if err != nil {
		return nil, err
	}

	// 3. Create first-boot marker
	marker := &schema_boot.FirstBootMarker{
		MachineID:     bm.Identity.MachineID,
		SchemaVersion: schema_system.CurrentVersion,
		GoldenHash:    goldenHash,
		Initialized:   true,
		CreatedAt:     time.Now(),
		BootTrust:     schema_system.TrustStrong,
	}
	if err := bm.Vault.SaveFirstBootMarker(marker); err != nil {
		return nil, err
	}

	// 4. Platform-specific auth
	authMgr := &auth.AuthManager{
		Vault:    bm.Vault,
		Platform: fullProfile.Platform.Final,
	}

	// Credential container
	type credential struct {
		UserID   string
		Password string
	}

	var cred credential

	found, err := bm.Vault.Read("credentials", bm.Identity.MachineID, &cred)
	if err != nil {
		return nil, fmt.Errorf("vault read failed: %w", err)
	}

	if !found {
		cred.UserID = bm.Identity.MachineID

		cred.Password, err = security_persistence.GenerateSecureKeyBase64()
		if err != nil {
			return nil, err
		}

		if err := bm.Vault.Write("credentials", cred.UserID, cred); err != nil {
			return nil, err
		}
	}

	// Authenticate
	session, err := authMgr.LoginOrSignUpInteractive()
	if err != nil {
		return nil, fmt.Errorf("auth failed during cold boot: %w", err)
	}

	capSet := BuildCapabilitySet(env.Platform.Final, resolveTier(session.Entity), resolveServiceProfile(env.Platform.Final))

	caps, err := boot.DetectDeviceCapabilities(env, capSet)
	if err != nil {
		return nil, fmt.Errorf("device capability detection failed: %w", err)
	}

	env.Discovery.Capabilities = *caps

	return &schema_system.BootSequence{
		Env:          env,
		Mode:         schema_boot.BootCold,
		Attested:     true,
		Capabilities: capSet,
		UserSession:  session,
		Service:      resolveServiceProfile(env.Platform.Final).Name,
		Tier:         resolveTier(session.Entity).Name,
		Entity:       session.Entity,
	}, nil
}

// ------------------------------------------------------------
// Fast Boot: use cached environment
// ------------------------------------------------------------
func (bm *boot_phase.BootManager) runFastBoot(env *schema_system.EnvConfig) (*schema_system.BootSequence, error) {
	logging.Info("[runFastBoot] Platform: %s | OS: %s | Arch: %s | EntityType: %v",
		env.Platform.Final, env.Identity.OS, env.Identity.Arch, bm.Identity.EntityType)

	marker, err := bm.Vault.LoadFirstBootMarker()
	if err != nil || marker.SchemaVersion != schema_system.CurrentVersion {
		return bm.runColdBoot(bm.Identity)
	}
	if err := security_verification.VerifyAgainstGolden(bm.Vault, marker.MachineID); err != nil {
		return bm.runColdBoot(bm.Identity)
	}

	raw, err := probe.IdentityProbe()
	if err != nil || raw.Identity.MachineID != env.Identity.MachineID || raw.Identity.OS != env.Identity.OS {
		return bm.runColdBoot(bm.Identity)
	}

	authMgr := &auth.AuthManager{Vault: bm.Vault, Platform: env.Platform.Final}

	var cred struct{ UserID, Password string }
	found, err := bm.Vault.Read("credentials", bm.Identity.MachineID, &cred)
	if !found || err != nil {
		return bm.runColdBoot(bm.Identity)
	}

	session, err := authMgr.LoginOrSignUpInteractive()
	if err != nil {
		return bm.runColdBoot(bm.Identity)
	}

	env.Attestation.SessionToken = session.SessionID

	// Capability merging
	capSet := BuildCapabilitySet(env.Platform.Final, resolveTier(session.Entity), resolveServiceProfile(env.Platform.Final))
	caps, err := boot.DetectDeviceCapabilities(env, capSet)
	if err != nil {
		return bm.runColdBoot(bm.Identity)
	}
	env.Discovery.Capabilities = *caps

	return &schema_system.BootSequence{
		Env:          env,
		Mode:         schema_boot.BootFast,
		Attested:     true,
		Capabilities: capSet,
		UserSession:  session,
		Service:      resolveServiceProfile(env.Platform.Final).Name,
		Tier:         resolveTier(session.Entity).Name,
		Entity:       session.Entity,
	}, nil
}

// ------------------ Helpers ------------------

func resolveTier(entity schema_system.EntityType) *schema_identity.TierProfile {

	switch entity {
	case schema_system.EntityOrganization:
		return &schema_identity.TierProfile{Name: schema_identity.TierEnterprise}

	case schema_system.EntityTester:
		return &schema_identity.TierProfile{Name: schema_identity.TierTester}

	default:
		return &schema_identity.TierProfile{Name: schema_identity.TierPersonal}
	}
}

func resolveServiceProfile(platform schema_system.PlatformClass) *schema_identity.ServiceProfile {

	switch platform {

	case schema_system.PlatformMobile:
		return &schema_identity.ServiceProfile{Name: schema_identity.ServicePersonal}

	case schema_system.PlatformVehicle:
		return &schema_identity.ServiceProfile{Name: schema_identity.ServiceMobility}

	case schema_system.PlatformIndustrial:
		return &schema_identity.ServiceProfile{Name: schema_identity.ServiceIndustrial}

	case schema_system.PlatformComputer:
		return &schema_identity.ServiceProfile{Name: schema_identity.ServicePersonal}

	default:
		return &schema_identity.ServiceProfile{Name: schema_identity.ServiceUnknown}
	}
}

// capSet := BuildCapabilitySet(bootSeq.Env.Platform.Final, resolveTier(bootSeq.Entity), resolveServiceProfile(bootSeq.Env.Platform.Final))
// BuildCapabilitySet computes platform + tier + service capabilities
func BuildCapabilitySet(platform schema_system.PlatformClass, tier *schema_identity.TierProfile, service *schema_identity.ServiceProfile,
) schema_security.CapabilitySet {
	var caps schema_security.CapabilitySet

	// Platform capabilities
	switch platform {
	case schema_system.PlatformVehicle:
		caps |= schema_security.CapCANBus | schema_security.CapSecureEnclave
	case schema_system.PlatformIndustrial:
		caps |= schema_security.CapIndustrialIO | schema_security.CapNetwork
	case schema_system.PlatformComputer:
		caps |= schema_security.CapLocalStorage | schema_security.CapNetwork | schema_security.CapBiometric
	}

	// Tier capabilities
	if tier.Name == schema_identity.TierEnterprise {
		caps |= schema_security.CapPersistentCloudLink
	}

	if service.Name == schema_identity.ServiceSystem {
		caps |= schema_security.CapSafetyCritical
	}
	return caps
}

func BootTrustToString(t schema_system.BootTrust) string {
	switch t {
	case schema_system.TrustStrong:
		return "strong"
	case schema_system.TrustWeak:
		return "weak"
	default:
		return "invalid"
	}
}
