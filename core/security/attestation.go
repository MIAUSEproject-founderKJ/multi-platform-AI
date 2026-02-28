//MIAUSEproject-founderKJ/multi-platform-AI/core/security/attestation.go

package security

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type VaultStore interface {
	LoadConfig(key string) (*schema.EnvConfig, error)
	SaveConfig(key string, cfg *schema.EnvConfig) error

	LoadGoldenHash(machine string) (string, error)
	LoadFirstBootMarker() (*schema.FirstBootMarker, error)
	SaveFirstBootMarker(*schema.FirstBootMarker) error
}

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

func VerifyEnvironment(v VaultStore, machineID string) error {

	hash, err := MeasureSelf()
	if err != nil {
		return fmt.Errorf("binary_measurement_failed: %w", err)
	}

	expected, err := v.LoadGoldenHash(machineID)
	if err != nil {
		return errors.New("baseline_missing")
	}

	if subtle.ConstantTimeCompare(hash, expected) != 1 {
		return errors.New("binary_tamper_detected")
	}

	return nil
}

func VerifyAgainstGolden(v VaultStore, machineID string) error {
	currentHash, err := MeasureSelf()
	if err != nil {
		return fmt.Errorf("failed_to_measure_binary: %w", err)
	}

	expectedHash, err := v.LoadGoldenHash(machineID)
	if err != nil {
		return errors.New("golden_hash_missing")
	}

	if subtle.ConstantTimeCompare(currentHash, expectedHash) != 1 {
		return errors.New("BINARY_TAMPER_DETECTED")
	}

	logging.Info("[SECURITY] Binary integrity verified.")
	return nil
}

func ProvisionGolden(v VaultStore, machineID string) error {
	hash, err := MeasureSelf()
	if err != nil {
		return err
	}

	return v.SealGoldenHash(machineID, hash)
}