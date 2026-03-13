// boot/phase_discovery.go
package boot

import (
	"fmt"

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

	env, err := probe.PassiveDiscovery()
	if err != nil {
		return nil, fmt.Errorf("passive scan failed: %w", err)
	}

	return &DiscoveryResult{
		InstanceID:   env.Identity.MachineName,
		PlatformType: env.Platform,
		OS:           env.OS,
		Architecture: env.Architecture,
	}, nil
}
