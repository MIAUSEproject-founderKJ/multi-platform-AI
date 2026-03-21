// boot/resolve_execution_context.go
package boot

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func ResolveExecutionContext(
	bs *schema.BootSequence,
) (*schema.BootContext, error) {

	// ------------------------------------------------------------
	// 1. Hard Validation (Fail Fast)
	// ------------------------------------------------------------

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

	// ------------------------------------------------------------
	// 2. Normalize Core Inputs (Single Source of Truth)
	// ------------------------------------------------------------

	env := bs.Env
	session := bs.UserSession

	entity := bs.Entity
	tier := bs.Tier
	service := bs.Service
	caps := bs.Capabilities

	// Trust resolution (CRITICAL DESIGN DECISION)
	// Prefer environment attestation as root-of-trust
	var trust schema.TrustLevel

	switch env.Attestation.Level {
	case schema.TrustStrong:
		trust = schema.TrustSystem
	case schema.TrustWeak:
		trust = schema.TrustDevice
	default:
		trust = schema.TrustUntrusted
	}

	// ------------------------------------------------------------
	// 3. Derive Permissions (Policy Engine)
	// ------------------------------------------------------------

	perms := make(map[schema.Permission]bool)

	// --- Base permission (always granted)
	perms[schema.PermUser] = true

	// --- Entity-based permissions
	switch entity {
	case schema.EntityPersonal:
		// baseline user only

	case schema.EntityOrganization:
		perms[schema.PermDiagnostics] = true

	case schema.EntityTester:
		perms[schema.PermDiagnostics] = true
		perms[schema.PermConfigEdit] = true

	case schema.EntityStranger:
		// minimal permissions only
	}

	// --- Tier-based permissions
	switch tier {
	case schema.TierEnterprise:
		perms[schema.PermDiagnostics] = true
		perms[schema.PermConfigEdit] = true
	}

	// --- Capability-based permissions (hardware-aware)
	if caps.Has(schema.CapCANBus) || caps.Has(schema.CapIndustrialIO) {
		perms[schema.PermHardwareIO] = true
	}

	if caps.Has(schema.CapSafetyCritical) && trust >= schema.TrustAdmin {
		perms[schema.PermSafetyOverride] = true
	}

	// --- Trust-based escalation (high-risk controls)
	if trust >= schema.TrustAdmin {
		perms[schema.PermAdmin] = true
		perms[schema.PermSafetyOverride] = true
	}

	// --- Merge session permissions (last layer, additive only)
	for p, allowed := range session.Permissions {
		if allowed {
			perms[p] = true
		}
	}

	// ------------------------------------------------------------
	// 4. Construct BootContext
	// ------------------------------------------------------------

	ctx := &schema.BootContext{
		PlatformClass: env.Platform.Final,
		Capabilities:  caps,
		Service:       service,
		Entity:        entity,
		Tier:          tier,
		BootMode:      bs.Mode,
		Permissions:   perms,
		TrustLevel:    trust,
	}

	return ctx, nil
}
