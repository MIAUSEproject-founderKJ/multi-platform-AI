// core/security/verification/verification_engine.go
package security_verification

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"os"

	security_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

func MeasureSelf() ([]byte, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(exePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func VerifyEnvironment(v security_persistence.VaultStore, machineID string) error {

	hash, err := MeasureSelf()
	if err != nil {
		return fmt.Errorf("binary_measurement_failed: %w", err)
	}

	expected, err := v.LoadGoldenHash(machineID)
	if err != nil {
		return errors.New("baseline_missing")
	}

	expectedBytes := []byte(expected)

	if subtle.ConstantTimeCompare(hash, expectedBytes) != 1 {
		return errors.New("binary_tamper_detected")
	}

	return nil
}

func VerifyAgainstGolden(v security_persistence.VaultStore, machineID string) error {
	currentHash, err := MeasureSelf()
	if err != nil {
		return fmt.Errorf("failed_to_measure_binary: %w", err)
	}

	expected, err := v.LoadGoldenHash(machineID)
	if err != nil {
		return errors.New("baseline_missing")
	}

	expectedBytes := []byte(expected)

	if subtle.ConstantTimeCompare(currentHash, expectedBytes) != 1 {
		return errors.New("binary_tamper_detected")
	}
	logging.Info("[SECURITY] Binary integrity verified.")
	return nil
}

func ProvisionGolden(v security_persistence.VaultStore, machineID string) ([]byte, error) {
	hash, err := MeasureSelf()
	if err != nil {
		return nil, fmt.Errorf("failed_to_measure_binary: %w", err)
	}
	return hash, err
}
