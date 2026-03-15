// boot/boot_modes.go

package boot

import (
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// DecideBootPath determines whether to run fast or cold boot
func (bm *BootManager) DecideBootPath() (*schema.BootSequence, error) {
	// Load last known environment
	lastkey := security.LastKnownEnvKey(bm.Identity.MachineID)
	env, err := bm.Vault.LoadConfig(lastkey)
	if err != nil {
		return bm.runColdBoot()
	}

	if env.SchemaVersion < schema.CurrentVersion {
		env = schema.Migrate(env)
	}

	// Verify golden baseline
	if _, err := bm.Vault.LoadGoldenHash(bm.Identity.MachineID); err != nil {
		return bm.runColdBoot()
	}

	// Perform fast boot
	return bm.runFastBoot(env)
}

// ------------------------------------------------------------
// Cold Boot: full hardware discovery and provisioning
// ------------------------------------------------------------
func (bm *BootManager) runColdBoot() (*schema.BootSequence, error) {
	// 1. Active hardware discovery
	fullProfile, err := probe.ActiveDiscovery(&bm.Identity.Hardware)
	if err != nil {
		return nil, fmt.Errorf("hardware discovery failed: %w", err)
	}

	bm.Identity.BindHardware(fullProfile)

	// 2. Provision golden baseline
	goldenHash, err := security.ProvisionGolden(bm.Vault, bm.Identity.MachineID)
	if err != nil {
		return nil, err
	}

	// 3. Create first-boot marker
	marker := &schema.FirstBootMarker{
		MachineName:   bm.Identity.MachineID,
		SchemaVersion: schema.CurrentVersion,
		GoldenHash:    goldenHash,
		Initialized:   true,
		CreatedAt:     time.Now(),
		TrustLevel:    schema.TrustStrong,
	}
	if err := bm.Vault.SaveFirstBootMarker(marker); err != nil {
		return nil, err
	}

	// 4. Platform-specific auth
	authMgr := &auth.AuthManager{
		Platform: fullProfile.Platform.Final,
		Entity:   bm.Identity.EntityType,
		Tier:     bm.Identity.TierType,
	}
	session, err := authMgr.LoginOrSignUp()
	if err != nil {
		return nil, fmt.Errorf("auth failed during cold boot: %w", err)
	}

	// 6. Attach session token
	fullProfile.Attestation.SessionToken = session.SessionID
	fullProfile.Attestation.Valid = true
	fullProfile.Attestation.Level = string(schema.TrustStrong)
	err = bm.Vault.SaveConfig(
		security.LastKnownEnvKey(bm.Identity.MachineID),
		fullProfile,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save last known environment: %w", err)
	}

	serviceType := resolveServiceProfile(fullProfile.Platform.Final)
	tierType := resolveTier(authMgr.Entity)

	capSet := BuildCapabilitySet(
		fullProfile.Platform.Final,
		tierType,
		serviceType,
	)

	return &schema.BootSequence{
		Env:          fullProfile,
		Mode:         schema.BootCold,
		Attested:     true,
		Capabilities: capSet,
		Service:      serviceType.Name,
		Tier:         tierType.Name,
		Entity:       authMgr.Entity,
	}, nil

}

// ------------------------------------------------------------
// Fast Boot: use cached environment
// ------------------------------------------------------------
func (bm *BootManager) runFastBoot(env *schema.EnvConfig) (*schema.BootSequence, error) {
	// 1. Verify against golden
	marker, err := bm.Vault.LoadFirstBootMarker()
	if err != nil || marker.SchemaVersion != schema.CurrentVersion {
		return bm.runColdBoot()
	}
	if err := security.VerifyAgainstGolden(bm.Vault, marker.MachineName); err != nil {
		return bm.runColdBoot()
	}

	// 2. Passive sanity scan
	raw, err := probe.IdentityProbe()
	if err != nil || raw.Identity.MachineID != env.Identity.MachineID || raw.Identity.OS != env.Identity.OS {
		return bm.runColdBoot()
	}

	// 3. Silent login
	authMgr := &auth.AuthManager{
		Platform: env.Platform.Final,
		Entity:   bm.Identity.EntityType,
		Tier:     bm.Identity.TierType,
	}
	session, err := authMgr.LoginOrSignUp()
	if err != nil {
		return nil, fmt.Errorf("auth failed during fast boot: %w", err)
	}

	// 4. Assign dynamic service & tier
	serviceType := resolveServiceProfile(env.Platform.Final)
	tierType := resolveTier(authMgr.Entity)

	// 5. Build capabilities
	capSet := BuildCapabilitySet(env.Platform.Final, tierType, serviceType)

	permList := security.DerivePermissions(
		env.Platform.Final,
		authMgr.Entity,
		string(tierType.Name),
	)

	permMap := make(map[string]bool)
	for _, p := range permList {
		permMap[p] = true
	}

	session.Permissions = permMap
	env.Attestation.SessionToken = session.SessionID

	return &schema.BootSequence{
		Env:          env,
		Mode:         schema.BootFast,
		Attested:     true,
		Capabilities: capSet,
		Service:      serviceType.Name,
		Tier:         tierType.Name,
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

	case schema.PlatformVehicle:
		return &schema.ServiceProfile{Name: schema.ServiceMobility}

	case schema.PlatformIndustrial:
		return &schema.ServiceProfile{Name: schema.ServiceIndustrial}

	case schema.PlatformComputer, schema.PlatformLaptop:
		return &schema.ServiceProfile{Name: schema.ServicePersonal}

	default:
		return &schema.ServiceProfile{Name: schema.ServiceUnknown}
	}
}

// BuildCapabilitySet computes platform + tier + service capabilities
func BuildCapabilitySet(
	platform schema.PlatformClass,
	tier *schema.TierProfile,
	service *schema.ServiceProfile,
) schema.CapabilitySet {
	var caps schema.CapabilitySet

	// Platform capabilities
	switch platform {
	case schema.PlatformVehicle:
		caps |= schema.CapCANBus | schema.CapSecureEnclave
	case schema.PlatformIndustrial:
		caps |= schema.CapIndustrialIO | schema.CapNetwork
	case schema.PlatformComputer, schema.PlatformLaptop:
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
