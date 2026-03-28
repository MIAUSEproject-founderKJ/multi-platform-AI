// boot/phase_discovery.go
package boot

import (
	"fmt"
	"context"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type DiscoveryResult struct {
	InstanceID   string
	PlatformType schema.PlatformClass
	OS           string
	Architecture string
}

func PhaseDiscovery() (*DiscoveryResult, error) {

	ctx := context.Background()

	env, err := probe.PassiveDiscovery(ctx)
	if err != nil {
		return nil, fmt.Errorf("passive scan failed: %w", err)
	}

	return &DiscoveryResult{
		InstanceID:   env.Identity.MachineID,
		PlatformType: env.Platform.Final,
		OS:           env.Identity.OS,
		Architecture: env.Identity.Arch,
	}, nil
}
