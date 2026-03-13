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
	env, err := bm.Vault.LoadConfig(lastKnownEnvKey)
	if err != nil || env.SchemaVersion != schema.CurrentVersion {
		return bm.runColdBoot()
	}

	// Verify golden baseline
	if _, err := bm.Vault.LoadGoldenHash(bm.Identity.MachineName); err != nil {
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
	goldenHash, err := security.ProvisionGolden(bm.Vault, bm.Identity.MachineName)
	if err != nil {
		return nil, err
	}

	// 3. Create first-boot marker
	marker := &schema.FirstBootMarker{
		MachineName:   bm.Identity.MachineName,
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

	// 5. Assign dynamic service & tier
	serviceProfile := resolveServiceProfile(fullProfile.Platform.Final)
	tierProfile := resolveTier(authMgr.Entity)

	// 6. Attach session token
	fullProfile.Attestation.SessionToken = session.SessionID
	fullProfile.Attestation.Valid = true
	fullProfile.Attestation.Level = string(schema.TrustStrong)
	if err := bm.Vault.SaveConfig(firstBootMarkerKey, fullProfile); err != nil {
		return nil, fmt.Errorf("failed to save first boot marker: %w", err)
	}

	// 7. Build capabilities
	capSet := BuildCapabilitySet(fullProfile.Platform.Final, tierProfile.Name, serviceProfile.Name)

	return &schema.BootSequence{
		Env:          fullProfile,
		Mode:         schema.BootCold,
		Attested:     true,
		Capabilities: capSet,
		Service:      core.ServiceType(serviceProfile.Name),
		Tier:         core.TierType(tierProfile.Name),
		Entity:       authMgr.Entity,
	}, nil

	err := bm.Vault.SaveConfig(lastKnownEnvKey(bm.Identity.MachineID), fullProfile)
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
	raw, err := probe.PassiveDiscovery()
	if err != nil || raw.MachineID != bm.Identity.MachineName || raw.OS != env.Identity.OS {
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
	serviceProfile := resolveServiceProfile(env.Platform.Final)
	tierProfile := resolveTier(authMgr.Entity)

	// 5. Build capabilities
	capSet := BuildCapabilitySet(env.Platform.Final, tierProfile.Name, serviceProfile.Name)
	session.Permissions = security.DerivePermissions(env.Platform.Final, authMgr.Entity, tierProfile.Name)
	env.Attestation.SessionToken = session.SessionID

	return &schema.BootSequence{
		Env:          env,
		Mode:         schema.BootFast,
		Attested:     true,
		Capabilities: capSet,
		Service:      core.ServiceType(serviceProfile.Name),
		Tier:         core.TierType(tierProfile.Name),
		Entity:       authMgr.Entity,
	}, nil
}

// ------------------ Helpers ------------------

func resolveTier(entity schema.EntityType) *core.TierProfile {
	if entity == schema.EntityOrganization || entity == schema.EntityTester {
		return &core.TierProfile{Name: "Funder"}
	}
	return &core.TierProfile{Name: "Non-Funder"}
}

func resolveServiceProfile(platform schema.PlatformClass) *core.ServiceProfile {
	switch platform {
	case schema.PlatformVehicle:
		return &core.ServiceProfile{Name: "AutonomousMobility"}
	case schema.PlatformIndustrial:
		return &core.ServiceProfile{Name: "IndustrialControl"}
	case schema.PlatformComputer, schema.PlatformLaptop:
		return &core.ServiceProfile{Name: "ProductivityAI"}
	default:
		return &core.ServiceProfile{Name: "GenericRuntime"}
	}
}

// BuildCapabilitySet computes platform + tier + service capabilities
func BuildCapabilitySet(platform schema.PlatformClass, tierName, serviceName string) core.CapabilitySet {
	var caps core.CapabilitySet

	// Platform capabilities
	switch platform {
	case schema.PlatformVehicle:
		caps |= core.CapCANBus | core.CapSecureEnclave
	case schema.PlatformIndustrial:
		caps |= core.CapIndustrialIO | core.CapNetwork
	case schema.PlatformComputer, schema.PlatformLaptop:
		caps |= core.CapLocalStorage | core.CapNetwork | core.CapBiometric
	}

	// Tier capabilities
	if tierName == "Funder" {
		caps |= core.CapPersistentCloudLink
	}

	// Service capabilities
	if serviceName == "AutonomousMobility" {
		caps |= core.CapSafetyCritical
	}

	return caps
}
