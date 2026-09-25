//bootstrap/resolver/boot_fast_path.go

package bootstrap_resolver

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/probe"
	core_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/verification"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

// ------------------------------------------------------------
// Fast Boot: use cached environment
// ------------------------------------------------------------

func (bm *BootManager) runFastBoot(
	env *internal_environment.EnvConfig,
	marker *internal_boot.FirstBootMarker,
) (*internal_boot.BootSequence, error) {
	env = internal_environment.Migrate(env)
	if env.SchemaVersion != internal_environment.CurrentVersion {
		// Migrate couldn't bring it forward (unknown/newer version) — don't guess.
		return bm.runColdBoot()
	}
	if err := core_verification.VerifyAgainstGolden(bm.Vault, marker.MachineID); err != nil {
		return bm.runColdBoot()
	}
	raw, err := probe.IdentityProbe()
	if err != nil || raw.Identity.MachineID != env.Identity.MachineID || raw.Identity.OS != env.Identity.OS {
		return bm.runColdBoot()
	}
	return &internal_boot.BootSequence{Env: env, Mode: internal_boot.BootFast, Attested: true}, nil
}
