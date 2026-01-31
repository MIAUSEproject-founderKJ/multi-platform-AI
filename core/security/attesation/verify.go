//core/security/attesation/verify.go
package attestation

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"multi-platform-AI/core/platform/probe"
	"multi-platform-AI/internal/logging"
)

// VerifyEnvironment performs a cryptographic check of the running binary.
func VerifyEnvironment(id *probe.HardwareIdentity) error {
	logging.Info("[SECURITY] Initiating Platform Attestation for: %s", id.PlatformType)

	// 1. MEASURED BOOT: Calculate the hash of the current executable
	currentHash, err := calculateSelfHash()
	if err != nil {
		return fmt.Errorf("failed_to_measure_binary: %w", err)
	}

	// 2. HW-BINDING: Compare with the "Golden Hash" stored in the Secure Vault
	// If this is a first boot, we "Provision" the hash.
	// If not, we "Verify" it.
	if err := verifyIntegrity(currentHash, id); err != nil {
		return fmt.Errorf("integrity_violation: %w", err)
	}

	logging.Info("[SECURITY] Attestation Successful. Environment Trusted.")
	return nil
}

func calculateSelfHash() ([]byte, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	
	data, err := os.ReadFile(exePath)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	return hash[:], nil
}

func verifyIntegrity(actualHash []byte, id *probe.HardwareIdentity) error {
	// SIMULATION: In a real TPM environment, we would ask the TPM 
	// to unseal a secret that only reveals itself if the PCR registers 
	// (which store the binary hash) match the expected state.
	
	// For now, we simulate a check against a signed 'manifest' file.
	expectedHash := getManifestHash(id.InstanceID)
	
	if string(actualHash) != string(expectedHash) {
		return errors.New("BINARY_TAMPER_DETECTED")
	}
	
	return nil
}