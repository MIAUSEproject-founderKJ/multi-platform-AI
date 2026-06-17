//bootstrap/resolver/boot_cold_path.go

package bootstrap_resolver

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/probe"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

// ------------------------------------------------------------
// Cold Boot: full hardware discovery and provisioning
// ------------------------------------------------------------
func (bm *BootManager) runColdBoot() (*internal_environment.BootSequence, error) {
	// 1. Active hardware discovery
	env := &internal_environment.EnvConfig{
		Identity: *bm.Identity,
	}

	fullProfile, err := probe.ActiveDiscovery(env)
	if err != nil {
		return nil, fmt.Errorf("hardware discovery failed: %w", err)
	}

	bm.Identity.BindHardware(fullProfile)

	return &internal_environment.BootSequence{
		Env:      fullProfile,
		Mode:     internal_boot.BootCold,
		Attested: true,
	}, nil

}
