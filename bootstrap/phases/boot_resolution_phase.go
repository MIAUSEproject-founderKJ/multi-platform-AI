//bootstrap/phases/boot_resolution_phase.go

package bootstrap_phase

import (
	bootstrap_resolver "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/resolver"
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/internal_environment/boot"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/internal_environment/environment"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
)

type BootManager struct {
	Vault    verification_persistence.VaultStore
	Identity *internal_environment.MachineIdentity
}

func PhaseBootResolution(identity *internal_environment.MachineIdentity) (*internal_environment.BootSequence, error) {
	bm := &BootManager{
		Identity: identity,
	}

	logging.Info(
		"[func PhaseContext] Platform: %s | OS: %s | Arch: %s | EntityType: %v",
		bm.Identity.PlatformType,
		bm.Identity.OS,
		bm.Identity.Arch,
		bm.Identity.EntityType,
	)

	bs, err := bootstrap_resolver.DecideBootPath(bm)
	if err != nil {
		return nil, err
	}

	if bs.Mode == internal_boot.BootCold {

		marker := &internal_boot.FirstBootMarker{
			MachineID: identity.MachineID,
		}

		if err := bm.Vault.MarkFirstBoot(marker.MachineID); err != nil {
			return nil, err
		}
	}

	return bs, nil
}
