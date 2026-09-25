// bootstrap/resolver/boot_path_resolver.go
package bootstrap_resolver

import (
	"errors"
	"fmt"

	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/keys"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
)

func (bm *BootManager) DecideBootPath() (*internal_boot.BootSequence, error) {
	if bm == nil || bm.Vault == nil || bm.Identity == nil || bm.Identity.MachineID == "" {
		return nil, errors.New("invalid BootManager state: missing vault or machine identity")
	}

	marker, err := bm.Vault.LoadFirstBootMarker()

	switch {
	case err == nil && marker.Initialized:
		// Authoritative signal says this machine completed a trusted boot before.
		return bm.attemptFastBoot(marker)

	case err == nil && !marker.Initialized:
		// Marker exists but was never completed — likely a crash mid cold-boot.
		// Fail safe: re-run cold boot rather than trust a half-written record.
		return bm.runColdBoot()

	case errors.Is(err, verification_persistence.ErrNotFound):
		// No marker at all — genuinely first boot.
		return bm.runColdBoot()

	default:
		// Vault unreadable/corrupted — do not guess. This is a hard failure,
		// not an ambiguous "treat as first boot" case.
		return nil, fmt.Errorf("failed to load first boot marker: %w", err)
	}
}

// attemptFastBoot loads the cached environment now that the marker has
// already confirmed this machine was previously attested. If the config
// is missing despite an initialized marker, that's an inconsistent state —
// fail safe to cold boot rather than trusting a partial fast path.
func (bm *BootManager) attemptFastBoot(marker *internal_boot.FirstBootMarker) (*internal_boot.BootSequence, error) {
	lastkey := keys.LastKnownEnvKey(bm.Identity.MachineID)
	env, err := bm.Vault.LoadConfig(lastkey)
	if err != nil {
		return bm.runColdBoot()
	}
	return bm.runFastBoot(env, marker)
}
