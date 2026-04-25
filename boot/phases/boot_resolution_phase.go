//boot/phases/boot_resolution_phase.go

package boot_phase

import (
	security_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

type BootProfile struct {
	Type string // FirstBoot | FastBoot | RecoveryBoot
}

func PhaseBootResolution(v security_persistence.VaultStore, identity *schema_system.MachineIdentity) (*schema_system.BootSequence, error) {
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

	if bs.Mode == schema_boot.BootCold {

		marker := &schema_boot.FirstBootMarker{
			MachineID: identity.MachineID,
		}

		if err := v.MarkFirstBoot(marker.MachineID); err != nil {
			return nil, err
		}
	}

	return bs, nil
}

type BootManager struct {
	Vault    security_persistence.VaultStore
	Identity *schema_system.MachineIdentity
}
