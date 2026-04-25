// boot\resolver\boot_context_resolver.go
package resolver

import (
	"fmt"

	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

func ResolveBootContext(bs *schema_system.BootSequence) (*schema_boot.BootContext, error) {

	if bs == nil {
		return nil, fmt.Errorf("boot sequence is nil")
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
	var trust schema_system.BootTrust

	switch env.Attestation.Level {
	case schema_system.TrustStrong:
		trust = schema_system.TrustStrong
	case schema_system.TrustWeak:
		trust = schema_system.TrustWeak
	default:
		trust = schema_system.TrustInvalid
	}

	// --- Permissions ---
	perms := make(map[schema_identity.Permission]bool)
	perms[schema_identity.PermUser] = true

	switch entity {
	case schema_system.EntityOrganization:
		perms[schema_identity.PermDiagnostics] = true
	case schema_system.EntityTester:
		perms[schema_identity.PermDiagnostics] = true
		perms[schema_identity.PermConfigEdit] = true
	}

	switch tier {
	case schema_identity.TierEnterprise:
		perms[schema_identity.PermDiagnostics] = true
		perms[schema_identity.PermConfigEdit] = true
	}

	if caps.Has(schema_security.CapCANBus) || caps.Has(schema_security.CapIndustrialIO) {
		perms[schema_identity.PermHardwareIO] = true
	}

	if caps.Has(schema_security.CapSafetyCritical) && trust == schema_system.TrustStrong {
		perms[schema_identity.PermSafetyOverride] = true
	}

	if trust == schema_system.TrustStrong {
		perms[schema_identity.PermAdmin] = true
		perms[schema_identity.PermSafetyOverride] = true
	}

	for p, allowed := range session.Permissions {
		if allowed {
			perms[p] = true
		}
	}

	ctx := &schema_boot.BootContext{
		PlatformClass: env.Platform.Final,
		Capabilities:  caps,
		Service:       service,
		Entity:        entity,
		Tier:          tier,
		BootMode:      bs.Mode,
		Permissions:   perms,
		TrustLevel:    schema_identity.TrustLevel(trust),
	}

	return ctx, nil
}
