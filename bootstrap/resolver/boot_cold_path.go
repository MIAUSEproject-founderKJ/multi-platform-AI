//bootstrap/resolver/boot_cold_path.go

package bootstrap_resolver

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/probe"
)

// ------------------------------------------------------------
// Cold Boot: full hardware discovery and provisioning
// ------------------------------------------------------------
func (bm *BootManager) runColdBoot() (*internal_environment.BootSequence, error) {
	// 1. Active hardware discovery
	fullProfile, err := probe.ActiveDiscovery(&bm.Identity.Hardware)
	if err != nil {
		return nil, fmt.Errorf("hardware discovery failed: %w", err)
	}

	bm.Identity.BindHardware(fullProfile)

	return &internal_environment.BootSequence{
		Env:      fullProfile,
		Mode:     internal_environment.BootCold,
		Attested: true,
	}, nil

}
