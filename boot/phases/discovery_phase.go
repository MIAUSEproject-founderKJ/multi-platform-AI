// boot\phases\discovery_phase.go
package boot

import (
	"context"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
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

	// passive discovery collects platform info, machine ID, OS, arch
	env, err := probe.PassiveDiscovery(ctx)
	if err != nil {
		return nil, fmt.Errorf("passive scan failed: %w", err)
	}

	// Ensure Platform.Final is set
	if env.Platform.Final == "" {
		env.Platform.Final = schema.PlatformComputer // default fallback
		logging.Warn("[DISCOVERY] Platform not detected, defaulting to 'computer'")
	}

	logging.Info("[phase_discovery.go] Platform: %s", env.Platform.Final)
	return &DiscoveryResult{
		InstanceID:   env.Identity.MachineID,
		PlatformType: env.Platform.Final,
		OS:           env.Identity.OS,
		Architecture: env.Identity.Arch,
	}, nil
}
