//boot/phase_context.go

package boot

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
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

func BuildBootContext(bs *schema.BootSequence) (*schema.BootContext, error) {

	if bs.Env == nil || !bs.Env.Attestation.Valid {
		return nil, fmt.Errorf("environment attestation failed")
	}

	perms := map[schema.Permission]bool{
		schema.PermUser: true,
	}

	if bs.Env.Attestation.Level == schema.TrustStrong {
		perms[schema.PermAdmin] = true
	}

	ctx := &schema.BootContext{
		PlatformClass: bs.Env.Identity.PlatformType,
		Capabilities:  bs.Capabilities,
		Service:       bs.Service,
		Entity:        bs.Entity,
		Tier:          bs.Tier,
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

		marker := &schema.FirstBootMarker{
			MachineID: identity.MachineID,
		}

		if err := v.MarkFirstBoot(marker.MachineID); err != nil {
			return nil, err
		}
	}

	return bs, nil
}

type BootManager struct {
	Vault    security.VaultStore
	Identity *schema.MachineIdentity
}
