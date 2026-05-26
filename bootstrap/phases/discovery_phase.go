// bootstrap/phases/discovery_phase.go
package bootstrap_phase

import (
	"context"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/probe"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
)

type DiscoveryResult struct {
	InstanceID   string
	PlatformType internal_environment.PlatformClass
	OS           string
	Architecture string
}

func PhaseDiscovery() (*DiscoveryResult, error) {

	ctx := context.Background()

	// passive discovery collects platform info, machine ID, OS, arch
	env, err := probe.PassiveDiscovery(ctx)
	if err != nil {
		return nil, fmt.Errorf("passive scan failed: %w", err)
	}

	logging.Info("[phase_discovery.go] Platform: %s", env.Platform.Final)
	return &DiscoveryResult{
		InstanceID:   env.Identity.MachineID,
		PlatformType: env.Platform.Final,
		OS:           env.Identity.OS,
		Architecture: env.Identity.Arch,
	}, nil
}
