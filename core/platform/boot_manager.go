//core/platform/boot_manager.go

package platform

import (
	"multi-platform-AI/core/platform/probe"
	"multi-platform-AI/core/security"
	"multi-platform-AI/internal/logging"
)

type BootMode string

const (
	ColdBoot BootMode = "COLD" // Full hardware discovery
	FastBoot BootMode = "FAST" // Optimized state resume
)

type BootManager struct {
	Vault    *security.IsolatedVault
	Platform *BootSequence
}

// ManageBoot determines the path: First-Boot or Subsequent-Boot.
func (bm *BootManager) ManageBoot() (*BootSequence, error) {
	// 1. Check for the "FirstBootMarker" in the Vault
	isFirstBoot := bm.Vault.IsMissingMarker("FirstBootMarker")

	if isFirstBoot {
		return bm.runColdBoot()
	}

	return bm.runFastBoot()
}

// runColdBoot (Stage 1): Full hardware discovery and user registration.
func (bm *BootManager) runColdBoot() (*BootSequence, error) {
	logging.Info("[BOOT] Stage 1 (Cold Boot): Initiating Full Discovery...")

	// Probing hardware (Lidar, CAN-bus, etc.)
	hwProfile := probe.AggressiveScan()

	// User Registration (Personal | Organization | Stranger)
	// If no IdentityToken, trigger Onboarding UI flow
	identity := bm.IdentifyUser()

	// Generate Security Token and Marker
	bm.Vault.WriteMarker("FirstBootMarker")
	bm.Vault.StoreToken("IdentityToken", identity.Token)

	return &BootSequence{Mode: "Discovery", PlatformID: hwProfile.ID}, nil
}

// runFastBoot (Stage 2): Optimized startup using persisted states.
func (bm *BootManager) runFastBoot() (*BootSequence, error) {
	logging.Info("[BOOT] Stage 2 (Fast Boot): Rapid Deployment...")

	// 1. Load Persisted EnvConfig (Skip full scan)
	config := bm.Vault.LoadConfig("LastKnownEnv")

	// 2. Hardware Heartbeat (Delta-Check)
	// Ensures sensors from the last session are still alive/unobstructed.
	if err := probe.Heartbeat(config); err != nil {
		logging.Warn("Hardware Delta detected! Reverting to Stage 1 Recovery.")
		return bm.runColdBoot()
	}

	// 3. Silent Login / Biometric Gate
	userRole := bm.ClassifyUserMatrix()

	return &BootSequence{
		Mode:       "Matured",
		UserRole:   userRole,
		TrustScore: config.LastTrustScore,
	}, nil
}

// ClassifyUserMatrix handles the "Transient" identity logic.
func (bm *BootManager) ClassifyUserMatrix() string {
	// Logic based on your matrix:
	// - Personal: Silent login via previous Biometric binding.
	// - Stranger: PIN entry / Guest mode (Robotaxi).
	// - Tester: Aggressive diagnostic handshake (Maintenance).
	return "OWNER" 
}