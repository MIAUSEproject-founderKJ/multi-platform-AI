// MIAUSEproject-founderKJ/multi-platform-AI/core/platform/boot.go

package platform

import (
	"errors"
	"fmt"
	"multi-platform-AI/core/platform/classify"
	"multi-platform-AI/core/platform/degrade"
	"multi-platform-AI/core/platform/probe"
	"multi-platform-AI/core/security/attestation"
	"multi-platform-AI/internal/logging"
)

// BootSequence defines the state of the platform during initialization
type BootSequence struct {
	PlatformID string
	TrustScore float64
	IsVerified bool
	Mode       string // "Autonomous", "Guarded", or "Discovery"
}

// RunBootSequence executes the PROBE -> CLASSIFY -> ATTEST logic.
func RunBootSequence() (*BootSequence, error) {
	logging.Info("Phase 1: Passive Hardware Probing...")
	
	// 1. PROBE: Passive scan (DMI, CPUID, Sysfs)
	rawHardware, err := probe.PassiveScan()
	if err != nil {
		return nil, fmt.Errorf("probe failure: %w", err)
	}

	logging.Info("Phase 2: Platform Classification...")
	
	// 2. CLASSIFY: Match hardware fingerprints to YAML types (Desktop/Vehicle/Industrial)
	identity, err := classify.Identify(rawHardware)
	if err != nil {
		// If identity is unknown, we don't crash; we enter Discovery Mode.
		logging.Warn("Unknown platform detected. Defaulting to Discovery Mode.")
		return &BootSequence{Mode: "Discovery", TrustScore: 0.1}, nil
	}

	logging.Info("Phase 3: Security Attestation & Vault Unlock...")

	// 3. ATTEST: Cryptographic check (TPM/TEE) to ensure software integrity
	if err := attestation.VerifyEnvironment(identity); err != nil {
		// If attestation fails (e.g., someone tampered with the code), force Safe Mode.
		return nil, degrade.ToSafeMode(errors.New("attestation_failed: integrity compromised"))
	}

	// 4. LOAD TRUST PRIOR: Check how reliable this specific instance was last time
	lastTrust := LoadTrustPrior(identity.InstanceID)

	return &BootSequence{
		PlatformID: identity.PlatformType,
		TrustScore: lastTrust,
		IsVerified: true,
		Mode:       DetermineExecutionMode(lastTrust),
	}, nil
}

// DetermineExecutionMode implements the Hysteresis logic for autonomy
func DetermineExecutionMode(trust float64) string {
	if trust > 0.95 {
		return "Autonomous"
	} else if trust > 0.50 {
		return "Guarded"
	}
	return "Discovery"
}