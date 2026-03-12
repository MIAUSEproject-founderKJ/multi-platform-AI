//boot/resolve_execution_context.go

package runtime

import (
	"fmt"


	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot"
)

func ResolveExecutionContext(
	bs *schema.BootSequence,
) (*boot.RuntimeContext, error) {

	// ------------------------------------------------------------
	// 1. Validate Machine Attestation
	// ------------------------------------------------------------

	if bs == nil || bs.Env == nil {
		return nil, fmt.Errorf("invalid boot sequence")
	}

	if !bs.Attested || !bs.Env.Attestation.Valid {
		return nil, fmt.Errorf("environment attestation invalid")
	}

	if bs.UserSession == nil {
		return nil, fmt.Errorf("missing authenticated session")
	}

	// ------------------------------------------------------------
	// 2. Trust Boot-Derived Values (Single Source of Truth)
	// ------------------------------------------------------------

	caps := bs.Capabilities
	entity := bs.Entity
	tier := bs.Tier
	service := bs.Service

	// ------------------------------------------------------------
	// 3. Derive Permissions (Policy Layer)
	// ------------------------------------------------------------

	var perms boot.PermissionSet

	// Base permission
	perms |= type.PermUser

	switch entity {
	case core.EntityAdmin:
		perms |= type.PermAdmin
	case core.EntityDevice:
		perms |= core.PermDeviceControl
	}

	if tier == core.TierEnterprise {
		perms |= core.PermFleetControl
	}

	// Optionally merge session-derived permissions
	perms |= bs.UserSession.Permissions

	// ------------------------------------------------------------
	// 4. Construct RuntimeContext
	// ------------------------------------------------------------

	return &boot.RuntimeContext{
		PlatformClass: bs.Env.Platform.Final,
		Capabilities:  caps,
		Service:       service,
		Entity:        entity,
		Tier:          tier,
		BootMode:      bs.Mode,
		Permissions:   perms,
	}, nil
}