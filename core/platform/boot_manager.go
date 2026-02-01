//core/platform/boot_manager.go

package platform

import (
	"fmt"
	"multi-platform-AI/api/hmi"
	"multi-platform-AI/core/policy"
	"multi-platform-AI/internal/logging"
	"multi-platform-AI/internal/mathutil"
	"multi-platform-AI/internal/monitor"
	"multi-platform-AI/plugins/navigation"
	"multi-platform-AI/plugins/perception"
)

// runColdBoot: Full Discovery with Active HMI & Vision Feedback
func (bm *BootManager) runColdBoot() (*BootSequence, error) {
	logging.Info("[BOOT] Stage 1 (Cold Boot): Full Discovery & Registration")

	// 1. Initialize Vision & SLAM (The "Eyes")
	vision := perception.NewVisionStream()
	nav := navigation.SLAMContext{}
	
	// 2. Start Vitals Monitoring (The "Nerves")
	vitals := monitor.NewVitalsMonitor(bm.Identity)
	vitals.Start()

	// 3. Link UI Pipelines (Connect Eyes, Nerves, and HUD)
	bm.attachVitalsToHMI(vitals.Stream)
	bm.linkVisionHUD(vision, vitals, &nav)

	// 4. Initialize SLAM using the live vision stream
	nav.InitializeSLAM(bm.Vault.Config, vision)

	// 5. Aggressive Probing (Progress: 10% - 50%)
	bm.updateProgress(0.1, "Waking up sensors (Lidar/CAN)...")
	fullProfile := probe.AggressiveScan(bm.Identity)

	bm.updateProgress(0.5, "Verifying Security Integrity...")
	// Perform actual attestation here
	if err := security.VerifyEnvironment(bm.Identity); err != nil {
    return nil, err
}

	// 6. User Onboarding (Progress: 80%)
	bm.updateProgress(0.8, "Awaiting User Identification...")
	userSession := bm.IdentifyUser()

	// 7. Bayesian Trust Decision
	evaluator := policy.TrustEvaluator{MinThreshold: mathutil.Q16FromFloat(0.9)}
	finalTrust := evaluator.Evaluate(bm.Vault.Config)

	// 8. Persistence & Final HMI Pulse
	bm.Vault.WriteMarker("FirstBootMarker")
	bm.Vault.StoreToken("IdentityToken", userSession.Token)
	bm.Vault.SaveConfig("LastKnownEnv", fullProfile)

	bm.updateProgress(1.0, fmt.Sprintf("Boot Complete. Trust: %s", determineLabel(finalTrust)))

	return &BootSequence{
		PlatformID: bm.Identity.PlatformType,
		TrustScore: finalTrust.Float64(),
		IsVerified: true,
		Mode:       "Discovery",
		UserRole:   userSession.Role,
	}, nil
}

// linkVisionHUD connects the internal data streams to the Perception Overlay
func (bm *BootManager) linkVisionHUD(vision *perception.VisionStream, vitals *monitor.VitalsMonitor, nav *navigation.SLAMContext) {
	// Link Vitals to HUD
	go func() {
		for v := range vitals.Stream {
			vision.UpdateVitals(v)
		}
	}()

	// Link Boot Progress to HUD
	go func() {
		// We use a local listener to the HMI pipe
		for p := range bm.HMIPipe {
			vision.UpdateProgress(p)
		}
	}()

	// Feed spatial markers (SLAM) back to HUD
	go func() {
		for {
    markers := nav.GetLatestMarkers()
    vision.UpdateSpatialMarkers(markers)
    time.Sleep(33 * time.Millisecond) // ADD THIS: Prevent 100% CPU usage (approx 30fps)
}
	}()
}

// ManageBoot determines the boot path based on the presence of the FirstBootMarker.
func (bm *BootManager) ManageBoot() (*BootSequence, error) {
    logging.Info("[BOOT_MGR] Resolving boot strategy for platform: %s", bm.Identity.PlatformType)

    // Check the vault for the FirstBootMarker
    isFirstBoot := bm.Vault.IsMissingMarker("FirstBootMarker")

    if isFirstBoot {
        // Path A: The system has never seen this hardware before or was reset.
        return bm.runColdBoot()
    }

    // Path B: The hardware is recognized. Attempt a high-speed resume.
    return bm.runFastBoot()
}


func (bm *BootManager) runFastBoot() (*BootSequence, error) {
    logging.Info("[BOOT] Stage 1 (Fast Boot): Resuming from Vault...")

    // 1. Quick Attestation
    // Even in fast boot, we must verify the binary hasn't been tampered with.
    if err := security.VerifyEnvironment(bm.Identity); err != nil {
        logging.Error("[BOOT] Security mismatch during resume. Redirecting to Cold Boot.")
        return bm.runColdBoot() 
    }

    // 2. Load Last Known Environment
    lastConfig, err := bm.Vault.LoadConfig("LastKnownEnv")
    if err != nil {
        return nil, fmt.Errorf("failed to load environment cache: %w", err)
    }

    // 3. Fast-track Bayesian Trust
    evaluator := policy.TrustEvaluator{MinThreshold: mathutil.Q16FromFloat(0.9)}
    finalTrust := evaluator.Evaluate(lastConfig)

    logging.Info("[BOOT] Fast Resume Complete. Trust Score: %.2f", finalTrust.Float64())

    return &BootSequence{
        PlatformID: bm.Identity.PlatformType,
        TrustScore: finalTrust.Float64(),
        IsVerified: true,
        Mode:       "Autonomous", // Assume autonomous on resume if trust is high
        UserRole:   "Operator",   // Pull from cached session if applicable
    }, nil
}

func (bm *BootManager) updateProgress(pct float32, msg string) {
	if bm.HMIPipe != nil {
		bm.HMIPipe <- hmi.ProgressUpdate{
			Percentage: pct,
			Message:    msg,
			Stage:      "BOOT_SEQUENCE",
		}
	}
}

func determineLabel(trust mathutil.Q16) string {
	val := trust.Float64()
	if val >= 0.9 { return "NOMINAL (High Confidence)" }
	if val >= 0.7 { return "DEGRADED (Caution)" }
	return "UNTRUSTED (Halt Required)"
}