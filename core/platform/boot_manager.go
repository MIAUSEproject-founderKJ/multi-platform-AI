//core/platform/boot_manager.go

package platform

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/policy"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/monitor"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/plugins/navigation"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/plugins/perception"
)

const CurrentSchemaVersion = 1

// ManageBoot handles the logic gate between Fast and Cold paths.
func (bm *BootManager) ManageBoot() (*schema.BootSequence, error) {
	// 1. Physical/Security Attestation
	if err := security.VerifyEnvironment(bm.Identity); err != nil {
		logging.Error("CRITICAL: Binary Integrity Compromised!")
		return nil, bm.EnterRecoveryMode(err) // Force safe state
	}

	isFirstBoot := bm.Vault.IsMissingMarker("FirstBootMarker")
	cachedEnv, err := bm.Vault.LoadConfig("LastKnownEnv")

	// 2. State Decision: If version jump or missing, re-probe.
	isOutdated := err == nil && cachedEnv.SchemaVersion != CurrentSchemaVersion
	if isFirstBoot || err != nil || isOutdated {
		if isOutdated {
			logging.Info("[BOOT] Schema jump (%d -> %d). Re-probing...", 
				cachedEnv.SchemaVersion, CurrentSchemaVersion)
		}
		return bm.runColdBoot()
	}

	return bm.runFastBoot()
}

func (bm *BootManager) runColdBoot() (*schema.BootSequence, error) {
	logging.Info("[BOOT] Stage 1 (Cold Boot): Full Discovery & Registration")

	// 1. Initial Subsystem Setup
	vision := perception.NewVisionStream()
	nav := navigation.SLAMContext{}
	vitals := monitor.NewVitalsMonitor(bm.Identity)
	vitals.Start()
	bm.linkVisionHUD(vision, vitals, &nav)

	// 2. Layered Hardware Probe
	bm.updateProgress(0.1, "Waking up sensors (Lidar/CAN)...")
	fullProfile, err := probe.ActiveDiscovery(bm.Identity.RawPassport)
	if err != nil {
		return nil, fmt.Errorf("[FATAL] Hardware detection failed: %w", err)
	}
	fullProfile.SchemaVersion = CurrentSchemaVersion

	// 3. Identity Finalization (The Scoring Engine)
	RunResolution(fullProfile)
	bm.Identity.Finalize(fullProfile)

	// 4. Bayesian Trust & Capability Gating
	evaluator := policy.TrustEvaluator{MinThreshold: 0.9}
	trustDesc := evaluator.Evaluate(fullProfile)

	// CAPABILITY GATING: Enforce hardware constraints on software modes
	if !fullProfile.Discovery.Capabilities.SupportsGoalControl {
		logging.Warn("[POLICY] Hardware lacks GoalControl. Restricting to MANUAL.")
		trustDesc.OperationMode = "MANUAL_ONLY"
	}
	if fullProfile.Discovery.Capabilities.SensorOnly {
		trustDesc.OperationMode = "READ_ONLY"
	}

	// 5. Persistence
	bm.Vault.WriteMarker("FirstBootMarker")
	bm.Vault.SaveConfig("LastKnownEnv", fullProfile)

	return &schema.BootSequence{
		PlatformID: string(bm.Identity.Config.Platform.Final),
		TrustScore: trustDesc.CurrentScore,
		IsVerified: true,
		Mode:       trustDesc.OperationMode,
		EnvConfig:  fullProfile,
	}, nil
}

func (bm *BootManager) runFastBoot() (*schema.BootSequence, error) {
	logging.Info("[BOOT] Stage 1 (Fast Boot): Resuming from Vault...")

	// Fast path skips hardware pings, relies on cached Attestation
	lastConfig, err := bm.Vault.LoadConfig("LastKnownEnv")
	if err != nil {
		return bm.runColdBoot() // Fallback if cache is corrupted
	}

	evaluator := policy.TrustEvaluator{MinThreshold: 0.9}
	trustDesc := evaluator.Evaluate(lastConfig)

	return &schema.BootSequence{
		PlatformID: string(bm.Identity.Config.Platform.Final),
		TrustScore: trustDesc.CurrentScore,
		IsVerified: true,
		Mode:       trustDesc.OperationMode,
		EnvConfig:  lastConfig,
	}, nil
}