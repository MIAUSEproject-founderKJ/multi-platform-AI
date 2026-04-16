// MIAUSEproject-founderKJ/multi-platform-AI/core/security/vault.go
// This file handles the low-level disk I/O for the secure markers.
package security

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/apppath"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type VaultStore interface {
	LoadConfig(key string) (*schema.EnvConfig, error)
	SaveConfig(key string, cfg *schema.EnvConfig) error
	MarkFirstBoot(machineID string) error
	LoadGoldenHash(machine string) (string, error)
	SealGoldenHash(machine string, hash []byte) error

	LoadFirstBootMarker() (*schema.FirstBootMarker, error)
	SaveFirstBootMarker(*schema.FirstBootMarker) error

	Read(key, id string, out interface{}) (bool, error)
	Write(key, id string, value interface{}) error
	Exists(key, id string) (bool, error)
}

type IsolatedVault struct {
	BaseDir string
	Key     []byte // Reserved for AES-GCM (Hardware-bound)
}

func GenerateSecureKeyBase64() (string, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(key), nil
}

// LoadSecureKey fetches the encryption key from environment variables.
// It ensures the key meets the 32-byte requirement for AES-256.
func LoadSecureKey() []byte {
	keyStr := os.Getenv("APP_ENCRYPTION_KEY")

	if keyStr == "" {
		// Auto-generate (dev or first boot)
		gen, err := GenerateSecureKeyBase64()
		if err != nil {
			log.Fatal("Failed to generate encryption key:", err)
		}

		log.Println("[SECURITY] Generated ephemeral encryption key")

		keyBytes, _ := base64.StdEncoding.DecodeString(gen)
		return keyBytes
	}

	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		log.Fatal("Invalid base64 key")
	}

	if len(key) != 32 {
		log.Fatalf("Key must decode to 32 bytes, got %d", len(key))
	}

	return key
}

// OpenVault initializes the secure directory with restricted owner-only access.
func OpenVault() (*IsolatedVault, error) {
	path := apppath.GetVaultPath()

	// Ensure 0700: Restricted to the user running the Kernel
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, fmt.Errorf("vault: failed to create directory: %w", err)
	}

	logging.Info("[VAULT] Secure storage initialized at %s", path)

	// Initialize your key at the start of the application
	encryptionKey := LoadSecureKey()

	fmt.Println("Success: Key loaded and validated.")

	return &IsolatedVault{
		BaseDir: path,
		Key:     encryptionKey,
	}, nil
}

func (v *IsolatedVault) SealGoldenHash(machine string, hash []byte) error {
	path := filepath.Join(v.BaseDir, "golden-"+machine)

	return os.WriteFile(path, hash, 0600)
}

func (v *IsolatedVault) LoadGoldenHash(machine string) (string, error) {
	path := filepath.Join(v.BaseDir, "golden-"+machine)

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DerivePermissions(
	platform schema.PlatformClass,
	entity schema.EntityType,
	tier schema.TierType,
) []schema.Permission {

	return []schema.Permission{schema.PermBasicRuntime}
}

// --- Config Logic ---

func (v *IsolatedVault) SaveConfig(name string, config *schema.EnvConfig) error {
	path := filepath.Join(v.BaseDir, name+".json")

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

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
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (v *IsolatedVault) StoreToken(name string, token string) error {
	path := filepath.Join(v.BaseDir, name)
	return os.WriteFile(path, []byte(token), 0600)
}

// --- Marker Logic ---

const firstBootVaultKey = "machine_first_boot_marker"

func (v *IsolatedVault) IsMissingMarker(name string) bool {
	_, err := os.Stat(filepath.Join(v.BaseDir, name))
	return os.IsNotExist(err)
}

func (v *IsolatedVault) WriteMarker(name string) error {
	path := filepath.Join(v.BaseDir, name)
	logging.Info("[VAULT] Sealing state marker: %s", name)
	return os.WriteFile(path, []byte("PROVISIONED"), 0600)
}

func (v *IsolatedVault) MarkFirstBoot(machineID string) error {

	marker := &schema.FirstBootMarker{
		MachineID:   machineID,
		Initialized: true,
		CreatedAt:   time.Now(),
	}

	return v.SaveFirstBootMarker(marker)
}

func (v *IsolatedVault) LoadFirstBootMarker() (*schema.FirstBootMarker, error) {

	path := filepath.Join(v.BaseDir, firstBootVaultKey+".json")

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load first boot marker: %w", err)
	}

	var marker schema.FirstBootMarker
	if err := json.Unmarshal(raw, &marker); err != nil {
		return nil, fmt.Errorf("failed to unmarshal first boot marker: %w", err)
	}

	return &marker, nil
}

func (v *IsolatedVault) SaveFirstBootMarker(marker *schema.FirstBootMarker) error {

	path := filepath.Join(v.BaseDir, firstBootVaultKey+".json")

	data, err := json.MarshalIndent(marker, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize first boot marker: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}
func (v *IsolatedVault) Read(collection, key string, out interface{}) (bool, error) {
	path := filepath.Join(v.BaseDir, collection+"_"+key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return false, err
	}
	return true, nil
}

func (v *IsolatedVault) Write(collection, key string, value interface{}) error {
	path := filepath.Join(v.BaseDir, collection+"_"+key+".json")
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (v *IsolatedVault) Exists(collection, key string) (bool, error) {
	path := filepath.Join(v.BaseDir, collection+"_"+key+".json")
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
