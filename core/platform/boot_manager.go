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

	// STEP 1: Authoritative user-first-boot check
	if bm.Vault.IsMissingMarker(firstBootMarker) {
		logging.Info("[BOOT] First launch detected")
		return bm.runColdBoot()
	}

	// STEP 2: Load cached config
	cachedEnv, err := bm.Vault.LoadConfig(lastKnownEnvKey)
	if err != nil {
		logging.Warn("[BOOT] Cached config missing or corrupted. Re-provisioning.")
		return bm.runColdBoot()
	}

	// STEP 3: Schema compatibility check
	if cachedEnv.SchemaVersion != currentSchemaVersion {
		logging.Info("[BOOT] Schema version mismatch. Re-provisioning.")
		return bm.runColdBoot()
	}

	// STEP 4: Fast path (no active hardware scan)
	return bm.runFastBoot(cachedEnv)
}

func (bm *BootManager) runColdBoot() (*schema.BootSequence, error) {

	logging.Info("[BOOT] Cold Boot: Full discovery + identity binding")

	fullProfile, err := probe.ActiveDiscovery(bm.Identity.RawPassport)
	if err != nil {
		return nil, fmt.Errorf("hardware discovery failed: %w", err)
	}

	fullProfile.SchemaVersion = currentSchemaVersion

	bm.Identity.Finalize(fullProfile)

	evaluator := policy.TrustEvaluator{MinThreshold: 0.9}
	trustDesc := evaluator.Evaluate(fullProfile)

	bm.Vault.WriteMarker(firstBootMarker)
	bm.Vault.SaveConfig(lastKnownEnvKey, fullProfile)

	return &schema.BootSequence{
		EnvConfig:  fullProfile,
		Mode:       trustDesc.OperationMode,
		Identity:   bm.Identity,
		TrustScore: trustDesc.CurrentScore,
		IsVerified: true,
	}, nil
}

func (bm *BootManager) runFastBoot(env *schema.EnvConfig) (*schema.BootSequence, error) {

	logging.Info("[BOOT] Fast Boot: Cached resume")

	// Optional: lightweight integrity check only
	if err := probe.SanityCheck(env); err != nil {
		logging.Warn("[BOOT] Hardware drift detected. Re-provisioning.")
		return bm.runColdBoot()
	}

	evaluator := policy.TrustEvaluator{MinThreshold: 0.9}
	trustDesc := evaluator.Evaluate(env)

	return &schema.BootSequence{
		EnvConfig:  env,
		Mode:       trustDesc.OperationMode,
		Identity:   bm.Identity,
		TrustScore: trustDesc.CurrentScore,
		IsVerified: true,
	}, nil
}