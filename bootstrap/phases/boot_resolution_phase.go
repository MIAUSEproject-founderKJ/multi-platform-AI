//bootstrap/phases/boot_resolution_phase.go

package bootstrap_phase

import (
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
)

func PhaseBootResolution(v verification_persistence.VaultStore, identity *internal_environment.MachineIdentity) (*internal_environment.BootSequence, error) {
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

	if bs.Mode == internal_boot.BootCold {

		marker := &internal_boot.FirstBootMarker{
			MachineID: identity.MachineID,
		}

		if err := v.MarkFirstBoot(marker.MachineID); err != nil {
			return nil, err
		}
	}

	return bs, nil
}

type BootManager struct {
	Vault    verification_persistence.VaultStore
	Identity *internal_environment.MachineIdentity
}
