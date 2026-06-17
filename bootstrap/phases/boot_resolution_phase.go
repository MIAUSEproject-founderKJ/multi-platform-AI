//bootstrap/phases/boot_resolution_phase.go

package bootstrap_phase

import (
	bootstrap_resolver "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/resolver"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/internal_environment/environment"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
)

// PhaseBootResolution determines the appropriate boot sequence based on the machine's identity.
func PhaseBootResolution(identity *internal_environment.MachineIdentity) (*internal_environment.BootSequence, error) {
	bm := &bootstrap_resolver.BootManager{
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

	return bs, nil
}
