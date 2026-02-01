//core/platform/boot_manager.go

package platform

import (
	"fmt"
	"time"

	"multi-platform-AI/api/hmi"
	"multi-platform-AI/core/policy"
	"multi-platform-AI/core/security"
	"multi-platform-AI/internal/logging"
	"multi-platform-AI/internal/monitor"
	"multi-platform-AI/plugins/navigation"
	"multi-platform-AI/plugins/perception"
	"multi-platform-AI/core/platform/probe" 
)

const CurrentSchemaVersion = 1

func (bm *BootManager) runColdBoot() (*BootSequence, error) {
	logging.Info("[BOOT] Stage 1 (Cold Boot): Full Discovery & Registration")

	// 1-3. Standard Setup
	vision := perception.NewVisionStream()
	nav := navigation.SLAMContext{}
	vitals := monitor.NewVitalsMonitor(bm.Identity)
	vitals.Start()
	bm.linkVisionHUD(vision, vitals, &nav)

	// 4. Hardware Probe
	bm.updateProgress(0.1, "Waking up sensors (Lidar/CAN)...")
	fullProfile, err := probe.AggressiveScan(bm.Identity)
	if err != nil {
		return nil, fmt.Errorf("[FATAL] Hardware detection failed: %w", err)
	}
	fullProfile.SchemaVersion = CurrentSchemaVersion

	// 5. Identity Finalization
	bm.updateProgress(0.4, "Finalizing Platform Reality...")
	bm.Identity.Finalize(fullProfile) 

	// 6. Security & User
	bm.updateProgress(0.6, "Verifying Security Integrity...")
	if err := security.VerifyEnvironment(bm.Identity); err != nil {
		return nil, err
	}
	userSession := bm.IdentifyUser()

	// 7. Bayesian Trust Decision (Unified)
	evaluator := policy.TrustEvaluator{MinThreshold: 0.9}
	trustDesc := evaluator.Evaluate(fullProfile) // Use the fresh profile

	logging.Info("[POLICY] Trust Score: %.2f | Mode: %s", trustDesc.CurrentScore, trustDesc.OperationMode)

	// 8. Persistence
	bm.Vault.WriteMarker("FirstBootMarker")
	bm.Vault.StoreToken("IdentityToken", userSession.Token)
	if err := bm.Vault.SaveConfig("LastKnownEnv", fullProfile); err != nil {
		logging.Warn("[BOOT] Failed to persist environment: %v", err)
	}

	bm.updateProgress(1.0, fmt.Sprintf("Boot Complete. Trust: %s", trustDesc.Label))

	return &BootSequence{
		PlatformID: bm.Identity.PlatformType,
		TrustScore: trustDesc.CurrentScore,
		IsVerified: true,
		Mode:       trustDesc.OperationMode, // Now returns "MANUAL_ONLY" or "AUTONOMOUS"
		UserRole:   userSession.Role,
	}, nil
}

func (bm *BootManager) runFastBoot() (*BootSequence, error) {
	logging.Info("[BOOT] Stage 1 (Fast Boot): Resuming from Vault...")

	// 1. Quick Attestation
	if err := security.VerifyEnvironment(bm.Identity); err != nil {
		logging.Error("[BOOT] Security mismatch. Redirecting to Cold Boot.")
		return bm.runColdBoot() 
	}

	// 2. Load cached profile
	lastConfig, err := bm.Vault.LoadConfig("LastKnownEnv")
	if err != nil {
		return nil, fmt.Errorf("failed to load environment cache: %w", err)
	}

	// 3. Bayesian Evaluation
	evaluator := policy.TrustEvaluator{MinThreshold: 0.9}
	trustDesc := evaluator.Evaluate(lastConfig)

	logging.Info("[BOOT] Fast Resume. Trust: %.2f | Mode: %s", trustDesc.CurrentScore, trustDesc.OperationMode)

	return &BootSequence{
		PlatformID: bm.Identity.PlatformType,
		TrustScore: trustDesc.CurrentScore,
		IsVerified: true,
		Mode:       trustDesc.OperationMode,
		UserRole:   "Operator", 
	}, nil
}

// ManageBoot refined with Schema Version Check from Reference
func (bm *BootManager) ManageBoot() (*BootSequence, error) {
	isFirstBoot := bm.Vault.IsMissingMarker("FirstBootMarker")
	
	// Attempt to peek at existing config to check version
	cachedEnv, err := bm.Vault.LoadConfig("LastKnownEnv")
	
	// Reference Logic: If version mismatch, force a Cold Boot (Re-probe)
	isOutdated := err == nil && cachedEnv.SchemaVersion != CurrentSchemaVersion

	if isFirstBoot || err != nil || isOutdated {
		if isOutdated {
			logging.Info("[BOOT] Schema mismatch. Triggering hardware re-probe...")
		}
		return bm.runColdBoot()
	}

	if err := security.VerifyEnvironment(bm.Identity); err != nil {
    logging.Error("CRITICAL: Binary Integrity Compromised!")
    // Force a recovery mode or halt execution
    return bm.EnterRecoveryMode(err) 
}

	return bm.runFastBoot()
}