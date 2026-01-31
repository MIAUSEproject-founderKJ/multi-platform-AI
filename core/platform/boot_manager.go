//core/platform/boot_manager.go

package platform

import (
	"multi-platform-AI/core/platform/probe"
	"multi-platform-AI/core/security"
	"multi-platform-AI/internal/logging"
)

type BootManager struct {
	Vault    *security.IsolatedVault
	Identity *probe.HardwareIdentity // Passed from the passive scan in boot.go
}

// ManageBoot handles the logic transition based on the FirstBootMarker.
func (bm *BootManager) ManageBoot() (*BootSequence, error) {
	// 1. Instant check for marker
	isFirstBoot := bm.Vault.IsMissingMarker("FirstBootMarker")

	if isFirstBoot {
		return bm.runColdBoot()
	}

	return bm.runFastBoot()
}

// runColdBoot (Stage 1): High-latency discovery for new installations.
func (bm *BootManager) runColdBoot() (*BootSequence, error) {
	logging.Info("[BOOT] Stage 1 (Cold Boot): Full Discovery & Registration")

	// 1. Aggressive Probing: Activate Lidar, CAN-bus handshake, etc.
	// We use the base identity to know which specific drivers to wake up.
	fullProfile := probe.AggressiveScan(bm.Identity)

	// 2. User Onboarding: Mandatory registration for first-ever boot.
	userSession := bm.IdentifyUser()

	// 3. Vault Seal: Save the state so Stage 2 is available next time.
	bm.Vault.WriteMarker("FirstBootMarker")
	bm.Vault.StoreToken("IdentityToken", userSession.Token)
	bm.Vault.SaveConfig("LastKnownEnv", fullProfile)

	return &BootSequence{
		PlatformID: bm.Identity.PlatformType,
		TrustScore: 0.1, // Initial discovery trust is low
		IsVerified: true,
		Mode:       "Discovery",
		UserRole:   userSession.Role,
	}, nil
}

// runFastBoot (Stage 2): Low-latency resumption for known environments.
func (bm *BootManager) runFastBoot() (*BootSequence, error) {
	logging.Info("[BOOT] Stage 2 (Fast Boot): Resuming Persisted State")

	// 1. Load context from Vault
	config := bm.Vault.LoadConfig("LastKnownEnv")

	// 2. Hardware Heartbeat: Delta-check against passive scan.
	// Does the current passive scan match the persisted AggressiveScan profile?
	if err := probe.Heartbeat(bm.Identity, config); err != nil {
		logging.Warn("Hardware mismatch detected (Portable move?). Reverting to Cold Boot.")
		return bm.runColdBoot()
	}

	// 3. Silent Login & User Classification
	userRole := bm.ClassifyUserMatrix()

	return &BootSequence{
		PlatformID: bm.Identity.PlatformType,
		TrustScore: config.LastTrustScore,
		IsVerified: true,
		Mode:       DetermineExecutionMode(config.LastTrustScore),
		UserRole:   userRole,
	}, nil
}

func (bm *BootManager) ClassifyUserMatrix() string {
	// Implementation of your Personal/Stranger/Tester logic
	return "OWNER"
}