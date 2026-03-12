//boot/phase_discovery.go
package boot

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type DiscoveryResult struct {
	Identity schema.MachineIdentity
	Platform schema.PlatformResolution
}

func PhaseDiscovery() (*DiscoveryResult, error) {

	env, err := probe.PassiveDiscovery()
	if err != nil {
		return nil, fmt.Errorf("passive scan failed: %w", err)
	}

	return &DiscoveryResult{
		Identity: env.Identity,
		Platform: env.Platform,
	}, nil
}