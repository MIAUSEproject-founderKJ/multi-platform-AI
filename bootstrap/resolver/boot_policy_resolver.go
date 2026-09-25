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

// ResolveExecutionContext is the single, canonical constructor for
// ExecutionContext. It is the only function permitted to call
// runtime_types.NewExecutionContext.
func ResolveExecutionContext(bs *internal_boot.BootSequence) (runtime_types.ExecutionContext, error) {
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
	entity, tier, service, caps := bs.Entity, bs.Tier, bs.Service, bs.Capabilities

	trust := deriveTrust(env.Attestation.Level)
	perms := derivePermissions(entity, tier, caps, trust, session)

	return runtime_types.NewExecutionContext(
		env.Platform.Final,
		caps,
		mapToTrustLevel(trust), // see note below on the two trust enums
		service,
		perms,
	), nil
}

func deriveTrust(level internal_common.BootTrust) internal_common.BootTrust {
	switch level {
	case internal_common.TrustStrong, internal_common.TrustWeak:
		return level
	default:
		return internal_common.TrustInvalid
	}
}

func derivePermissions(
	entity internal_common.EntityKind,
	tier internal_common.TierType,
	caps internal_environment.CapabilitySet,
	trust internal_common.BootTrust,
	session *user_setting.UserSession,
) map[user_setting.PermissionKey]bool {
	perms := map[user_setting.PermissionKey]bool{user_setting.PermUser: true}

	switch entity {
	case internal_common.EntityOrganization:
		perms[user_setting.PermDiagnostics] = true
	case internal_common.EntityTester:
		perms[user_setting.PermDiagnostics] = true
		perms[user_setting.PermConfigEdit] = true
	}

	if tier == internal_common.TierEnterprise {
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

	// Session permissions can only narrow, never grant beyond derived policy —
	// this is the "who never overrides where" rule from the Constitution.
	// The original code's `if allowed && perms[p]` already enforces this
	// correctly (session can't add permissions policy didn't derive); keep it.
	sessionPerms := session.Claims.Permissions

	for p := range perms {
		if !sessionPerms[p] {
			delete(perms, p)
		}
	}
	return perms
}
