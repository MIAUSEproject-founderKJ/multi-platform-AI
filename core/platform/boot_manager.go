//core/platform/boot_manager.go

package platform

import (
	"errors"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type BootManager struct {
	Vault    VaultStore
	Identity *schema.Identity
}

func (bm *BootManager) DecideBootPath() (*schema.BootSequence, error) {

	// Attempt to load previous environment (fast boot)
	env, err := bm.Vault.LoadConfig(lastKnownEnvKey)
	if err != nil || env.SchemaVersion != currentSchemaVersion {
		return bm.runColdBoot()
	}

	// Golden baseline check
	if _, err := bm.Vault.LoadGoldenHash(bm.Identity.MachineName); err != nil {
		return bm.runColdBoot()
	}

	return bm.runFastBoot(env)
}

// ------------------------------------------------------------
// Cold Boot: full discovery, first-time setup, registration
// ------------------------------------------------------------
func (bm *BootManager) runColdBoot() (*schema.BootSequence, error) {

	// 1. Active hardware discovery
	fullProfile, err := probe.ActiveDiscovery(bm.Identity.MachineName)
	if err != nil {
		return nil, fmt.Errorf("hardware discovery failed: %w", err)
	}
	fullProfile.SchemaVersion = currentSchemaVersion
	bm.Identity.BindHardware(fullProfile)

	// 2. Provision golden baseline
	if err := security.ProvisionGolden(bm.Vault, bm.Identity.MachineName); err != nil {
		return nil, err
	}

	// 3. Save discovered environment
	if err := bm.Vault.SaveConfig(lastKnownEnvKey, fullProfile); err != nil {
		return nil, err
	}

	// 4. Platform-specific login/sign-up
	auth := &AuthManager{
		Platform: fullProfile.Platform.Final,
		Entity:   bm.Identity.EntityType,
		Tier:     bm.Identity.TierType,
	}
	session, err := auth.LoginOrSignUp()
	if err != nil {
		return nil, fmt.Errorf("auth failed during cold boot: %w", err)
	}

	// 5. Assign dynamic service & tier profiles
	serviceProfile := resolveServiceProfile(fullProfile.Platform.Final)
	tierProfile := resolveTier(auth.Entity)

	// 6. Build capability graph based on platform + tier + service
	capSet := capability.BuildCapabilitySet(fullProfile.Platform.Final, tierProfile.Name, serviceProfile.Name)

	// 7. Formalize first-boot marker
	fullProfile.Attestation.SessionToken = session.SessionID
	fullProfile.Attestation.Valid = true
	fullProfile.Attestation.Level = string(schema.TrustStrong)
	if err := bm.Vault.SaveConfig(firstBootMarkerKey, fullProfile); err != nil {
		return nil, fmt.Errorf("failed to save first boot marker: %w", err)
	}

	// 8. Construct BootSequence with all runtime context
	return &schema.BootSequence{
		Env:      fullProfile,
		Mode:     schema.BootCold,
		Attested: true,
	}, nil
}

func (bm *BootManager) runFastBoot(env *schema.EnvConfig) (*schema.BootSequence, error) {

	// Verify environment against golden baseline
	if err := security.VerifyAgainstGolden(bm.Vault, bm.Identity.MachineName); err != nil {
		return bm.runColdBoot()
	}

	if err := bm.sanityCheck(env); err != nil {
		return bm.runColdBoot()
	}

	// Silent login / auto-login
	auth := &AuthManager{
		Platform: env.Platform.Final,
		Entity:   bm.Identity.EntityType,
		Tier:     bm.Identity.TierType,
	}
	session, err := auth.LoginOrSignUp()
	if err != nil {
		return nil, fmt.Errorf("auth failed during fast boot: %w", err)
	}

	// Dynamic service and tier assignment
	serviceProfile := resolveServiceProfile(env.Platform.Final)
	tierProfile := resolveTier(auth.Entity)

	// Build capability graph
	capSet := capability.BuildCapabilitySet(env.Platform.Final, tierProfile.Name, serviceProfile.Name)

	// Update session and env with assigned capabilities
	session.Permissions = security.DerivePermissions(env.Platform.Final, auth.Entity, tierProfile.Name)
	env.Attestation.SessionToken = session.SessionID

	return &schema.BootSequence{
		Env:      env,
		Mode:     schema.BootFast,
		Attested: true,
	}, nil
}

func (bm *BootManager) sanityCheck(env *schema.EnvConfig) error {
	raw, err := probe.PassiveScan()
	if err != nil {
		return err
	}

	if raw.InstanceID != bm.Identity.MachineName {
		return errors.New("machine_identity_changed")
	}

	if raw.PlatformType != env.PlatformClass {
		return errors.New("platform_class_drift")
	}

	return nil
}

func (bm *BootManager) checkFirstBootMarker() (bool, error) {
	marker, err := bm.Vault.LoadFirstBootMarker()
	if errors.Is(err, security.ErrNotFound) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return !marker.Initialized, nil
}

func resolveTier(entity schema.EntityType) *core.TierProfile {
	switch entity {
	case schema.EntityOrganization, schema.EntityTester:
		return &core.TierProfile{Name: "Funder"}
	default:
		return &core.TierProfile{Name: "Non-Funder"}
	}
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

func BuildCapabilitySet(platform schema.PlatformClass, tierName, serviceName string) core.CapabilitySet {
	var caps core.CapabilitySet

	// Base platform capabilities
	switch platform {
	case schema.PlatformVehicle:
		caps |= core.CapCANBus | core.CapSecureEnclave
	case schema.PlatformIndustrial:
		caps |= core.CapIndustrialIO | core.CapNetwork
	case schema.PlatformComputer, schema.PlatformLaptop:
		caps |= core.CapLocalStorage | core.CapNetwork | core.CapBiometric
	}

	// Tier-based capabilities
	if tierName == "Funder" {
		caps |= core.CapPersistentCloudLink
	}

	// Service-specific (example)
	if serviceName == "AutonomousMobility" {
		caps |= core.CapSafetyCritical
	}

	return caps
}
