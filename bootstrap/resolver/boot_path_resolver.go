//bootstrap/resolver/boot_path_resolver.go

package bootstrap_resolver

import (
	"fmt"

	bootstrap_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/phases"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/keys"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

// DecideBootPath determines whether to run fast or cold boot
func (bm *bootstrap_phase.BootManager) DecideBootPath() (*internal_environment.BootSequence, error) {
	// Load last known environment
	lastkey := keys.LastKnownEnvKey(bm.Identity.MachineID)
	env, err := bm.Vault.LoadConfig(lastkey)
	if err != nil {
		return nil, fmt.Errorf("failed to load last known environment: %w", err)
	}

	if env == nil {
		return bm.runColdBoot()
	}
	// Perform fast boot
	return bm.runFastBoot(env)
}
