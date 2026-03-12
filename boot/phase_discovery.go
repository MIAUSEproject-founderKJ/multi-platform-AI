//boot/phase_discovery.go

package boot

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
)

type DiscoveryResult struct {
	InstanceID   string
	PlatformType string
	OS           string
	Architecture string
}

func PhaseDiscovery() (*DiscoveryResult, error) {

	raw, err := probe.PassiveScan()
	if err != nil {
		return nil, fmt.Errorf("passive scan failed: %w", err)
	}

	return &DiscoveryResult{
		InstanceID:   raw.InstanceID,
		PlatformType: raw.PlatformType,
		OS:           raw.OS,
		Architecture: raw.Architecture,
	}, nil
}