//boot/phase_context.go

package boot

import (
	"errors"
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/optimization"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/capability"
)
type Capability string

const (
	CapCANBus              Capability = "CAN_BUS"
	CapBiometric           Capability = "BIOMETRIC"
	CapHighFreqSensor      Capability = "HIGH_FREQ_SENSOR"
	CapFileSystem          Capability = "FILE_SYSTEM"
	CapMicrophone          Capability = "MICROPHONE"
	CapSafetyCritical      Capability = "SAFETY_CRITICAL"
	CapPersistentCloudLink Capability = "PERSISTENT_CLOUD"
)

type BootProfile struct {
	Type string // FirstBoot | FastBoot | RecoveryBoot
}

type RuntimeContext struct {
	PlatformClass schema.PlatformClass
	Capabilities  core.CapabilitySet
	Service       core.ServiceType
	Entity        core.EntityType
	Tier          core.TierType
	BootMode      core.BootMode
	Permissions   boot.PermissionSet
	Router        *router.Router
	Optimizer     optimization.Optimizer
}


func BuildRuntimeContext(bs *schema.BootSequence) (*RuntimeContext, error){

	if bs.Attestation == nil || !bs.Attestation.Valid {
		return nil, fmt.Errorf("environment attestation failed")
	}

	var caps core.CapabilitySet
	var perms boot.PermissionSet

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
	perms |= type.PermUser

	if bs.Attestation.Level == schema.TrustStrong {
		perms |= type.PermAdmin
	}

	ctx := &RuntimeContext{
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


func PhaseContext(v *security.IsolatedVault, identity *schema.MachineIdentity) (*schema.BootSequence, error) {

    bm := &BootManager{
        Vault:    v,
        Identity: identity,
    }

    bs, err := bm.DecideBootPath()
    if err != nil {
        return nil, err
    }

    if bs.Mode == schema.BootCold {
        if err := v.MarkFirstBoot(identity.MachineName); err != nil {
            return nil, err
        }
    }

    return bs, nil
}

type BootManager struct {
	Vault    VaultStore
	Identity *schema.Identity
}

