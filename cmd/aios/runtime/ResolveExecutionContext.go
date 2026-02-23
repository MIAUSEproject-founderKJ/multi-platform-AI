//cmd/aios/runtime/ResolveExecutionContext.go

package runtime

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func ResolveExecutionContext(
	bs *schema.BootSequence,
	vault *security.IsolatedVault,
) (*core.RuntimeContext, error) {

	if bs.Env.Attestation.Valid == false {
		return nil, fmt.Errorf("environment attestation invalid")
	}

	var caps core.CapabilitySet
	var perms core.PermissionSet

	// ------------------------------------------------------------------
	// 1. Derive Capabilities from Hardware
	// ------------------------------------------------------------------

	for _, bus := range bs.Env.Hardware.Buses {
		if bus.Type == "can" && bus.Confidence > 40000 {
			caps |= core.CapCANBus
		}
	}

	if bs.Env.Discovery.Capabilities.SupportsAcceleratedCompute {
		caps |= core.CapSecureEnclave
	}

	if bs.Env.Hardware.HasBattery {
		caps |= core.CapLocalStorage
	}

	// ------------------------------------------------------------------
	// 2. Resolve Entity (from vault/session token)
	// ------------------------------------------------------------------

	entity := core.EntityUser
	tier := core.TierFree

	meta, err := vault.LoadUserMetadata()
	if err == nil {
		entity = meta.Entity
		tier = meta.Tier
	}

	// ------------------------------------------------------------------
	// 3. Derive Permissions
	// ------------------------------------------------------------------

	perms |= core.PermUser

	if entity == core.EntityAdmin {
		perms |= core.PermAdmin
	}

	if entity == core.EntityDevice {
		perms |= core.PermDeviceControl
	}

	if tier == core.TierEnterprise {
		perms |= core.PermFleetControl
	}

	// ------------------------------------------------------------------
	// 4. Derive Service From Platform
	// ------------------------------------------------------------------

	service := core.ServicePersonal

	switch bs.Env.Platform.Final {
	case schema.PlatformVehicle,
		schema.PlatformDrone,
		schema.PlatformRobot:

		service = core.ServiceSystem

	case schema.PlatformIndustrial:
		service = core.ServiceEnterprise
	}

	return &core.RuntimeContext{
		PlatformClass: bs.Env.Platform.Final,
		Capabilities:  caps,
		Service:       service,
		Entity:        entity,
		Tier:          tier,
		BootMode:      bs.Mode,
		Permissions:   perms,
	}, nil
}
