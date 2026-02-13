// MIAUSEproject-founderKJ/multi-platform-AI/core/security/vault.go
// This file handles the low-level disk I/O for the secure markers.
package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/apppath"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type IsolatedVault struct {
	BaseDir string
	Key     []byte // Reserved for AES-GCM (Hardware-bound)
}

// OpenVault initializes the secure directory with restricted owner-only access.
func OpenVault() (*IsolatedVault, error) {
	path := apppath.GetVaultPath()

	// Ensure 0700: Restricted to the user running the Kernel
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, fmt.Errorf("vault: failed to create directory: %w", err)
	}

	logging.Info("[VAULT] Secure storage initialized at %s", path)

	return &IsolatedVault{
		BaseDir: path,
		// TODO: Implement Key Derivation (Argon2) from Hardware UUID
		Key: []byte("temporary-32-byte-dev-key-12345"),
	}, nil
}

// --- Marker Logic ---

func (v *IsolatedVault) IsMissingMarker(name string) bool {
	_, err := os.Stat(filepath.Join(v.BaseDir, name))
	return os.IsNotExist(err)
}

func (v *IsolatedVault) WriteMarker(name string) error {
	path := filepath.Join(v.BaseDir, name)
	logging.Info("[VAULT] Sealing state marker: %s", name)
	return os.WriteFile(path, []byte("PROVISIONED"), 0600)
}

// --- Config Logic ---

func (v *IsolatedVault) SaveConfig(name string, config *schema.EnvConfig) error {
	path := filepath.Join(v.BaseDir, name+".json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (v *IsolatedVault) LoadConfig(name string) (*schema.EnvConfig, error) {
	path := filepath.Join(v.BaseDir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config schema.EnvConfig
	err = json.Unmarshal(data, &config)
	return &config, err
}

func (v *IsolatedVault) StoreToken(name string, token string) error {
	path := filepath.Join(v.BaseDir, name)
	return os.WriteFile(path, []byte(token), 0600)
}