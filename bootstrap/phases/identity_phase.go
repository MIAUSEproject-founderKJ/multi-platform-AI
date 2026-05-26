// bootstrap/phases/identity_phase.go
package bootstrap_phase

import (
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
)

func PhaseIdentity(d *DiscoveryResult) (*internal_environment.MachineIdentity, error) {

	logging.Info("[phase_identity] Platform: %s", d.PlatformType)

	identity := &internal_environment.MachineIdentity{
		MachineID:    d.InstanceID,
		PlatformType: d.PlatformType,
		OS:           d.OS,
		Arch:         d.Architecture,
	}

	return identity, nil
}
