// boot\phases\identity_phase.go
package boot

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func PhaseIdentity(d *DiscoveryResult) (*schema.MachineIdentity, error) {

	logging.Info("[phase_identity] Platform: %s", d.PlatformType)

	identity := &schema.MachineIdentity{
		MachineID:    d.InstanceID,
		PlatformType: d.PlatformType,
		OS:           d.OS,
		Arch:         d.Architecture,
	}

	return identity, nil
}
