//boot/phase_context.go

package boot

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)



type BootProfile struct {
	Type string // FirstBoot | FastBoot | RecoveryBoot
}

func PhaseContext(v security.VaultStore, identity *schema.MachineIdentity) (*schema.BootSequence, error) {
	bm := &BootManager{
		Vault:    v,
		Identity: identity,
	}

	logging.Info(
		"[func PhaseContext] Platform: %s | OS: %s | Arch: %s | EntityType: %v",
		bm.Identity.PlatformType,
		bm.Identity.OS,
		bm.Identity.Arch,
		bm.Identity.EntityType,
	)

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
