//core/platform/boot_manager.go

package platform

import (
	"fmt"
	"multi-platform-AI/api/hmi" // Assuming Protobuf generated path
	"multi-platform-AI/core/platform/probe"
	"multi-platform-AI/core/security"
	"multi-platform-AI/internal/logging"
	"multi-platform-AI/internal/monitor"
)

type BootManager struct {
	Vault           *security.IsolatedVault
	Identity        *probe.HardwareIdentity
	HMIPipe         chan hmi.ProgressUpdate
	currentProgress float32
}

func (bm *BootManager) ManageBoot() (*BootSequence, error) {
	isFirstBoot := bm.Vault.IsMissingMarker("FirstBootMarker")

	if isFirstBoot {
		return bm.runColdBoot()
	}
	return bm.runFastBoot()
}

// runColdBoot: Full Discovery with Active HMI Feedback
func (bm *BootManager) runColdBoot() (*BootSequence, error) {
	logging.Info("[BOOT] Stage 1 (Cold Boot): Full Discovery & Registration")
vision := perception.NewVisionStream()
	// 1. Initialize Vitals & HMI Feedback Loop
	vitals := monitor.NewVitalsMonitor(bm.Identity)
	vitals.Start()
	bm.attachVitalsToHMI(vitals.Stream)
nav := navigation.SLAMContext{}
nav.InitializeSLAM(bm.Config, vision)
	// 2. Aggressive Probing (Progress: 10% - 50%)
	bm.updateProgress(0.1, "Waking up sensors (Lidar/CAN)...")
	fullProfile := probe.AggressiveScan(bm.Identity)
	
	bm.updateProgress(0.5, "Verifying Security Integrity...")
	// Security attestation logic happens here...

	// 3. User Onboarding (Progress: 80%)
	bm.updateProgress(0.8, "Awaiting User Identification...")
	userSession := bm.IdentifyUser()

	// 4. Persistence
	bm.Vault.WriteMarker("FirstBootMarker")
	bm.Vault.StoreToken("IdentityToken", userSession.Token)
	bm.Vault.SaveConfig("LastKnownEnv", fullProfile)

	bm.updateProgress(1.0, "Boot Complete. Welcome.")
	return &BootSequence{
		PlatformID: bm.Identity.PlatformType,
		TrustScore: 0.1,
		IsVerified: true,
		Mode:       "Discovery",
		UserRole:   userSession.Role,
	}, nil
}

// Internal helper to pipe monitor data to the HMI channel
func (bm *BootManager) attachVitalsToHMI(vitalsStream <-chan monitor.SystemVitals) {
	go func() {
		for v := range vitalsStream {
			bm.HMIPipe <- hmi.ProgressUpdate{
				Message:    fmt.Sprintf("CPU: %.1f%% | VRAM: %dMB | Temp: %.1f°C", v.CPULoad, v.VRAMUsage/1024/1024, v.Temperature),
				Percentage: bm.currentProgress,
			}
		}
	}()
}

func (bm *BootManager) updateProgress(p float32, msg string) {
	bm.currentProgress = p
	bm.HMIPipe <- hmi.ProgressUpdate{
		Percentage: p,
		Message:    msg,
	}
}

// Link Vitals to HUD
go func() {
    for v := range vitals.Stream {
        vision.UpdateVitals(v)
    }
}()

// Link Boot Progress to HUD
go func() {
    for p := range bm.HMIPipe {
        vision.UpdateProgress(p)
    }
}()

// Feed spatial markers back to HUD
go func() {
    for {
        marker := nav.GetLatestMarkers()
        vision.UpdateSpatialMarkers(marker)
    }
}()