// boot/boot_modes.go

package boot

import (
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// DecideBootPath determines whether to run fast or cold boot
func (bm *BootManager) DecideBootPath() (*schema.BootSequence, error) {
	// Load last known environment
	lastkey := security.LastKnownEnvKey(bm.Identity.MachineID)
	env, err := bm.Vault.LoadConfig(lastkey)
	if err != nil {
		return bm.runColdBoot(bm.Identity)
	}

	if env.SchemaVersion < schema.CurrentVersion {
		env = schema.Migrate(env)
		_ = bm.Vault.SaveConfig(lastkey, env)
	}

	// Verify golden baseline
	if _, err := bm.Vault.LoadGoldenHash(bm.Identity.MachineID); err != nil {
		return bm.runColdBoot(bm.Identity)
	}

	// Fast boot
	return bm.runFastBoot(env)
}

// ------------------------------------------------------------
// Cold Boot: full hardware discovery and provisioning
// ------------------------------------------------------------
func (bm *BootManager) runColdBoot(identity *schema.MachineIdentity) (*schema.BootSequence, error) {
	// Use the discovered platform + identity info
	env := &schema.EnvConfig{
		Platform: schema.PlatformResolution{
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
	goldenHash, err := security.ProvisionGolden(bm.Vault, bm.Identity.MachineID)
	if err != nil {
		return nil, err
	}

	// 3. Create first-boot marker
	marker := &schema.FirstBootMarker{
		MachineID:     bm.Identity.MachineID,
		SchemaVersion: schema.CurrentVersion,
		GoldenHash:    goldenHash,
		Initialized:   true,
		CreatedAt:     time.Now(),
		BootTrust:     schema.TrustStrong,
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

		cred.Password, err = security.GenerateSecureKeyBase64()
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
	
capSet := BuildCapabilitySet(env.Platform.Final, resolveTier(authMgr.Entity), resolveServiceProfile(env.Platform.Final))

caps, err := DetectDeviceCapabilities(env, capSet)
if err != nil {
    return nil, fmt.Errorf("device capability detection failed: %w", err)
}

env.Discovery.Capabilities = *caps

	return &schema.BootSequence{
		Env:         env,
		Mode:        schema.BootCold,
		Attested:    true,
		Capabilities: capSet,
		UserSession: session,
		Service:     resolveServiceProfile(env.Platform.Final).Name,
		Tier:        resolveTier(authMgr.Entity).Name,
		Entity:      authMgr.Entity,
	}, nil
}



// ------------------------------------------------------------
// Fast Boot: use cached environment
// ------------------------------------------------------------
func (bm *BootManager) runFastBoot(env *schema.EnvConfig) (*schema.BootSequence, error) {
    logging.Info("[runFastBoot] Platform: %s | OS: %s | Arch: %s | EntityType: %v",
        env.Platform.Final, env.Identity.OS, env.Identity.Arch, bm.Identity.EntityType)

    marker, err := bm.Vault.LoadFirstBootMarker()
    if err != nil || marker.SchemaVersion != schema.CurrentVersion {
        return bm.runColdBoot(bm.Identity)
    }
    if err := security.VerifyAgainstGolden(bm.Vault, marker.MachineID); err != nil {
        return bm.runColdBoot(bm.Identity)
    }

    raw, err := probe.IdentityProbe()
    if err != nil || raw.Identity.MachineID != env.Identity.MachineID || raw.Identity.OS != env.Identity.OS {
        return bm.runColdBoot(bm.Identity)
    }

    authMgr := &auth.AuthManager{Vault: bm.Vault, Platform: env.Platform.Final}

    var cred struct { UserID, Password string }
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
    capSet := BuildCapabilitySet(env.Platform.Final, resolveTier(authMgr.Entity), resolveServiceProfile(env.Platform.Final))
    caps, err := DetectDeviceCapabilities(env, capSet)
    if err != nil {
        return bm.runColdBoot(bm.Identity)
    }
    env.Discovery.Capabilities = *caps

    return &schema.BootSequence{
        Env:          env,
        Mode:         schema.BootFast,
        Attested:     true,
        Capabilities: capSet,
		UserSession: session,
        Service:      resolveServiceProfile(env.Platform.Final).Name,
        Tier:         resolveTier(authMgr.Entity).Name,
        Entity:       authMgr.Entity,
    }, nil
}
// ------------------ Helpers ------------------

func resolveTier(entity schema.EntityType) *schema.TierProfile {

	switch entity {
	case schema.EntityOrganization:
		return &schema.TierProfile{Name: schema.TierEnterprise}

	case schema.EntityTester:
		return &schema.TierProfile{Name: schema.TierTester}

	default:
		return &schema.TierProfile{Name: schema.TierPersonal}
	}
}

func resolveServiceProfile(platform schema.PlatformClass) *schema.ServiceProfile {

	switch platform {

	case schema.PlatformMobile:
		return &schema.ServiceProfile{Name: schema.ServicePersonal}

	case schema.PlatformVehicle:
		return &schema.ServiceProfile{Name: schema.ServiceMobility}

	case schema.PlatformIndustrial:
		return &schema.ServiceProfile{Name: schema.ServiceIndustrial}

	case schema.PlatformComputer:
		return &schema.ServiceProfile{Name: schema.ServicePersonal}

	default:
		return &schema.ServiceProfile{Name: schema.ServiceUnknown}
	}
}

//capSet := BuildCapabilitySet(bootSeq.Env.Platform.Final, resolveTier(bootSeq.Entity), resolveServiceProfile(bootSeq.Env.Platform.Final))
// BuildCapabilitySet computes platform + tier + service capabilities
func BuildCapabilitySet(platform schema.PlatformClass,tier *schema.TierProfile,service *schema.ServiceProfile,
) schema.CapabilitySet {
	var caps schema.CapabilitySet

	// Platform capabilities
	switch platform {
	case schema.PlatformVehicle:
		caps |= schema.CapCANBus | schema.CapSecureEnclave
	case schema.PlatformIndustrial:
		caps |= schema.CapIndustrialIO | schema.CapNetwork
	case schema.PlatformComputer:
		caps |= schema.CapLocalStorage | schema.CapNetwork | schema.CapBiometric
	}

	// Tier capabilities
	if tier.Name == schema.TierEnterprise {
		caps |= schema.CapPersistentCloudLink
	}

	if service.Name == schema.ServiceSystem {
		caps |= schema.CapSafetyCritical
	}
	return caps
}

func BootTrustToString(t schema.BootTrust) string {
	switch t {
	case schema.TrustStrong:
		return "strong"
	case schema.TrustWeak:
		return "weak"
	default:
		return "invalid"
	}
}
