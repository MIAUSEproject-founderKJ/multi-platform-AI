//bootstrap/resolver/boot_cold_path.go

package bootstrap_resolver

import (
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/probe"
	core_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/verification"
	keys "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/keys"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

// ------------------------------------------------------------
// Cold Boot: full hardware discovery and provisioning
// ------------------------------------------------------------

func (bm *BootManager) runColdBoot() (*internal_boot.BootSequence, error) {
	env := &internal_environment.EnvConfig{
		SchemaVersion: internal_environment.CurrentVersion,
		GeneratedAt:   time.Now(),
		Identity:      *bm.Identity,
	}

	fullProfile, err := probe.ActiveDiscovery(env)
	if err != nil {
		return nil, fmt.Errorf("hardware discovery failed: %w", err)
	}

	bm.Identity.BindHardware(fullProfile)

	// Measure and seal the golden baseline — this is what attestation checks against.
	hash, err := core_verification.MeasureSelf()
	if err != nil {
		return nil, fmt.Errorf("binary measurement failed: %w", err)
	}
	if err := bm.Vault.SealGoldenHash(bm.Identity.MachineID, hash); err != nil {
		return nil, fmt.Errorf("failed to seal golden hash: %w", err)
	}

	// Persist config so the next boot can take the fast path.
	if err := bm.Vault.SaveConfig(keys.LastKnownEnvKey(bm.Identity.MachineID), fullProfile); err != nil {
		return nil, fmt.Errorf("failed to persist environment config: %w", err)
	}

	// Complete the marker — this is the record that cold boot actually finished.
	marker := &internal_boot.FirstBootMarker{
		MachineID:     bm.Identity.MachineID,
		SchemaVersion: internal_environment.CurrentVersion,
		GoldenHash:    hash,
		Initialized:   true,
		CreatedAt:     time.Now(),
		BootTrust:     internal_common.TrustStrong, // or derive from discovery confidence
	}
	if err := bm.Vault.MarkFirstBoot(marker); err != nil {
		return nil, fmt.Errorf("failed to finalize first boot marker: %w", err)
	}

	return &internal_boot.BootSequence{
		Env:      fullProfile,
		Mode:     internal_boot.BootCold,
		Attested: true,
	}, nil
}
