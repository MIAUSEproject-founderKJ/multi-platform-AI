//core/platform/boot_manager.go

package platform

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/policy"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type BootManager struct {
	Vault    VaultStore
	Identity *schema.Identity
}

func (bm *BootManager) DecideBootPath() (*schema.BootSequence, error) {

	firstBoot, err := bm.checkFirstBootMarker()
	if err != nil {
		return nil, err
	}

	if firstBoot {
		return bm.runColdBoot()
	}

	// Non-first boot → perform measured verification
	err = security.VerifyEnvironment(bm.Vault, bm.Identity.MachineName)
	if errors.Is(err, ErrBaselineMissing) {
		// Marker exists but baseline missing = tamper state
		return bm.runRecoveryBoot()
	}
	if err != nil {
		return nil, fmt.Errorf("measured boot failure: %w", err)
	}

	return bm.runNormalBoot()
}

func (bm *BootManager) checkFirstBootMarker() (bool, error) {
	marker, err := bm.Vault.LoadFirstBootMarker()
	if errors.Is(err, security.ErrNotFound) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return !marker.Initialized, nil
}


func (bm *BootManager) runColdBoot() (*schema.BootSequence, error) {

	// Initialize baseline measurements
	if err := security.EstablishBaseline(bm.Vault, bm.Identity); err != nil {
		return nil, err
	}

	if err := bm.Vault.StoreFirstBootMarker(&schema.FirstBootMarker{
		Initialized: true,
	}); err != nil {
		return nil, err
	}

	return &schema.BootSequence{
		Mode:     schema.ModeCold,
		Identity: bm.Identity,
	}, nil
}

func (bm *BootManager) runNormalBoot() (*schema.BootSequence, error) {
	return &schema.BootSequence{
		Mode:     schema.ModeNormal,
		Identity: bm.Identity,
	}, nil
}

func (bm *BootManager) sanityCheck(env *schema.EnvConfig) error {
	raw, err := probe.PassiveScan()
	if err != nil {
		return err
	}

	if raw.InstanceID != bm.Identity.MachineName {
		return errors.New("machine_identity_changed")
	}

	if raw.PlatformType != env.PlatformClass {
		return errors.New("platform_class_drift")
	}

	return nil
}

func (bm *BootManager) runRecoveryBoot() (*schema.BootSequence, error) {
	return &schema.BootSequence{
		Mode:     schema.ModeRecovery,
		Identity: bm.Identity,
	}, nil
}