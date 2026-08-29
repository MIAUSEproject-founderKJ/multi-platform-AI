// bootstrap/resolver/boot_policy_resolver.go
package bootstrap_resolver

import (
	"fmt"

	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

func ResolveBootContext(bs *internal_boot.BootSequence) (*runtime_types.ExecutionContext, error) {

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
	var trust internal_boot.FirstBootMarker

	switch env.Attestation.Level {
	case internal_common.TrustStrong:
		trust = internal_common.TrustStrong
	case internal_common.TrustWeak:
		trust = internal_common.TrustWeak
	default:
		trust = internal_common.TrustInvalid
	}

	// --- Permissions ---
	perms := make(map[user_setting.PermissionKey]bool)
	perms[user_setting.PermUser] = true

	switch entity {
	case internal_common.EntityOrganization:
		perms[user_setting.PermDiagnostics] = true
	case internal_common.EntityTester:
		perms[user_setting.PermDiagnostics] = true
		perms[user_setting.PermConfigEdit] = true
	}

	switch tier {
	case internal_common.TierEnterprise:
		perms[user_setting.PermDiagnostics] = true
		perms[user_setting.PermConfigEdit] = true
	}

	if caps.Has(internal_environment.CapCANBus) || caps.Has(internal_environment.CapIndustrialIO) {
		perms[user_setting.PermHardwareIO] = true
	}

	if caps.Has(internal_environment.CapSafetyCritical) && trust == internal_common.TrustStrong {
		perms[user_setting.PermAdmin] = true
		perms[user_setting.PermSafetyOverride] = true
	}

	for p, allowed := range session.Permissions {
		if allowed && perms[p] {
			perms[p] = true // keep allowed
		}
	}

	bootctx := &bootstrap_resolver.ExecutionContext{}

	return bootctx, nil
}
