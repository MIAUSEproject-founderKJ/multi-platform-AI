// boot/phases/identity_phase.go
package boot_phase

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

func PhaseIdentity(d *DiscoveryResult) (*schema_system.MachineIdentity, error) {

	logging.Info("[phase_identity] Platform: %s", d.PlatformType)

	identity := &schema_system.MachineIdentity{
		MachineID:    d.InstanceID,
		PlatformType: d.PlatformType,
		OS:           d.OS,
		Arch:         d.Architecture,
	}

	return identity, nil
}
