//boot/bootstrap.go

package boot

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/schema"
)

func BuildRuntimeContext(bs *schema.BootSequence) (*core.RuntimeContext, error) {

	if bs.Attestation == nil || !bs.Attestation.Valid {
		return nil, fmt.Errorf("environment attestation failed")
	}

	var caps core.CapabilitySet
	var perms core.PermissionSet

	// Platform capability resolution
	switch bs.Identity.Platform {

	case schema.PlatformWorkstation:
		caps |= core.CapNetwork | core.CapLocalStorage | core.CapBiometric

	case schema.PlatformAutomotive:
		caps |= core.CapCANBus | core.CapSecureEnclave

	case schema.PlatformIndustrial:
		caps |= core.CapIndustrialIO | core.CapNetwork

	case schema.PlatformRobotics:
		caps |= core.CapIndustrialIO | core.CapSecureEnclave
	}

	// Default permission logic (example)
	perms |= core.PermUser

	if bs.Attestation.Level == schema.TrustStrong {
		perms |= core.PermAdmin
	}

	ctx := &core.RuntimeContext{
		PlatformClass: bs.Identity.Platform,
		Capabilities:  caps,
		Service:       core.ServicePersonal,
		Entity:        core.EntityUser,
		Tier:          core.TierFree,
		BootMode:      bs.Mode,
		Permissions:   perms,
	}

	return ctx, nil
}
