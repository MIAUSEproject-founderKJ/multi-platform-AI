// bootstrap/resolver/boot_policy_resolver.go
package bootstrap_resolver

import (
	"fmt"

	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

func ResolveBootContext(bs *internal_environment.BootSequence) (*bootstrap.BootContext, error) {

	if bs == nil {
		return nil, fmt.Errorf("bootstrap sequence is nil")
	}

	if bs.Env == nil {
		return nil, fmt.Errorf("missing environment config")
	}

	if !bs.Attested || !bs.Env.Attestation.Valid {
		return nil, fmt.Errorf("environment attestation invalid")
	}

	if bs.UserSession == nil {
		return nil, fmt.Errorf("missing authenticated session")
	}

	env := bs.Env
	session := bs.UserSession

	entity := bs.Entity
	tier := bs.Tier
	service := bs.Service
	caps := bs.Capabilities

	// --- Trust resolution ---
	var trust internal_environment.BootTrust

	switch env.Attestation.Level {
	case internal_environment.TrustStrong:
		trust = internal_environment.TrustStrong
	case internal_environment.TrustWeak:
		trust = internal_environment.TrustWeak
	default:
		trust = internal_environment.TrustInvalid
	}

	// --- Permissions ---
	perms := make(map[user_setting.PermissionKey]bool)
	perms[user_setting.PermUser] = true

	switch entity {
	case internal_environment.EntityOrganization:
		perms[user_setting.PermDiagnostics] = true
	case internal_environment.EntityTester:
		perms[user_setting.PermDiagnostics] = true
		perms[user_setting.PermConfigEdit] = true
	}

	switch tier {
	case user_setting.TierEnterprise:
		perms[user_setting.PermDiagnostics] = true
		perms[user_setting.PermConfigEdit] = true
	}

	if caps.Has(internal_environment.CapCANBus) || caps.Has(internal_environment.CapIndustrialIO) {
		perms[user_setting.PermHardwareIO] = true
	}

	if caps.Has(internal_environment.CapSafetyCritical) && trust == internal_environment.TrustStrong {
		perms[user_setting.PermSafetyOverride] = true
	}

	if trust == internal_environment.TrustStrong {
		perms[user_setting.PermAdmin] = true
		perms[user_setting.PermSafetyOverride] = true
	}

	for p, allowed := range session.Permissions {
		if allowed {
			perms[p] = true
		}
	}

	ctx := &bootstrap.BootContext{
		PlatformClass: env.Platform.Final,
		Capabilities:  caps,
		Service:       service,
		Entity:        entity,
		Tier:          tier,
		BootMode:      bs.Mode,
		Permissions:   perms,
		TrustLevel:    user_setting.TrustLevel(trust),
	}

	return ctx, nil
}
