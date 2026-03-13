// boot/phase_identity.go
package boot

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/i  nternal/schema"

type DiscoveryResult struct {
	InstanceID   string
	PlatformType schema.PlatformClass
	OS           string
	Architecture string
}

func PhaseIdentity(d *DiscoveryResult) (*schema.MachineIdentity, error) {

	identity := &schema.MachineIdentity{
		MachineName: d.InstanceID,
		Platform:    d.PlatformType,
		OS:          d.OS,
		Arch:        d.Architecture,
	}

	return identity, nil
}
