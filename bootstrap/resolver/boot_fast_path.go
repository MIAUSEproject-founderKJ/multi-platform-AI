//bootstrap/resolver/boot_fast_path.go

package bootstrap_resolver

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/probe"
	core_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/verification"
)

// ------------------------------------------------------------
// Fast Boot: use cached environment
// ------------------------------------------------------------
func (bm *BootManager) runFastBoot(env *internal_environment.EnvConfig) (*internal_environment.BootSequence, error) {
	// 1. Verify against golden
	marker, err := bm.Vault.LoadFirstBootMarker()
	if err != nil || internal_environment.Version != internal_environment.CurrentVersion {
		return bm.runColdBoot()
	}
	if err := core_verification.VerifyAgainstGolden(bm.Vault, marker.MachineID); err != nil {
		return bm.runColdBoot()
	}

	// 2. Passive sanity scan
	raw, err := probe.IdentityProbe()
	if err != nil || raw.Identity.MachineID != env.Identity.MachineID || raw.Identity.OS != env.Identity.OS {
		return bm.runColdBoot()
	}

	return &internal_environment.BootSequence{
		Env:      env,
		Mode:     internal_environment.BootFast,
		Attested: true,
	}, nil
}
