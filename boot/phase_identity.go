// boot/phase_identity.go
package boot

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

func PhaseIdentity(d *DiscoveryResult) (*schema.MachineIdentity, error) {

	identity := &schema.MachineIdentity{
		MachineID:    d.InstanceID,
		PlatformType: d.PlatformType,
		OS:           d.OS,
		Arch:         d.Architecture,
	}

	return identity, nil
}
