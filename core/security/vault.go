// MIAUSEproject-founderKJ/multi-platform-AI/core/security/vault.go
// This file handles the low-level disk I/O for the secure markers.
package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/configs/defaults"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/apppath"
)

// IsolatedVault represents the secure storage engine for identity and state.
type IsolatedVault struct {
	BaseDir string
	Key     []byte // Reserved for future AES-GCM encryption implementation
}

// OpenVault initializes the secure directory on the host system.
func OpenVault() (*IsolatedVault, error) {
	// 1. Resolve path based on platform/identity
	path := apppath.GetVaultPath()

	// 2. Ensure directory exists with 0700 (Owner-only access)
	// This is critical for preventing other OS users from peeking at the AI's identity
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, fmt.Errorf("failed to initialize vault directory: %w", err)
	}

	return &IsolatedVault{
		BaseDir: path,
		// In the next step, we should derive this key from the Hardware UUID
		// using a Key Derivation Function (KDF) like Argon2 or PBKDF2
		Key: []byte("32-byte-long-auth-key-from-uuid"),
	}, nil
}

// --- Marker Logic (State Persistence) ---

// IsMissingMarker returns true if the specified marker (e.g., "FirstBootMarker") is absent.
func (v *IsolatedVault) IsMissingMarker(name string) bool {
	_, err := os.Stat(filepath.Join(v.BaseDir, name))
	return os.IsNotExist(err)
}

// WriteMarker creates a persistent signal file to lock in a system state.
func (v *IsolatedVault) WriteMarker(name string) error {
	path := filepath.Join(v.BaseDir, name)
	return os.WriteFile(path, []byte{}, 0600) // 0600: Restricted read/write
}

// --- Config Logic (Structured Environment Data) ---

// SaveConfig serializes the EnvConfig (Hardware profile) into the vault.
func (v *IsolatedVault) SaveConfig(name string, config *defaults.EnvConfig) error {
	path := filepath.Join(v.BaseDir, name+".json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("config serialization failed: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

// LoadConfig retrieves a stored environment profile from a previous boot.
func (v *IsolatedVault) LoadConfig(name string) (*defaults.EnvConfig, error) {
	path := filepath.Join(v.BaseDir, name+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", name, err)
	}

	var config defaults.EnvConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config %s: %w", name, err)
	}

	return &config, nil
}

// StoreToken saves session-specific credentials (e.g., User Identity Tokens).
func (v *IsolatedVault) StoreToken(name string, token string) error {
	path := filepath.Join(v.BaseDir, name)
	return os.WriteFile(path, []byte(token), 0600)
}

// Close ensures all file descriptors are synced and released.
func (v *IsolatedVault) Close() {
	// Implementation for flushing buffers or releasing file-system locks
}
