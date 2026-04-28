// bootstrap/phases/identity_phase.go
package bootstrap_phase

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
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
