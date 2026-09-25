// core/security/persistence/vault_store.go
package verification_persistence

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/apppath"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	"go.uber.org/zap"
)

var (
	ErrNotFound      = errors.New("vault: record not found")
	ErrPathTraversal = errors.New("vault: invalid key path outside base directory")
)

type IsolatedVault struct {
	BaseDir string
	Key     []byte
}

func OpenVault(log *zap.Logger) (*IsolatedVault, error) {
	path := apppath.GetVaultPath()
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, fmt.Errorf("vault directory creation failed (%s): %w", path, err)
	}
	key, err := LoadSecureKey()
	if err != nil {
		return nil, fmt.Errorf("vault key initialization failed: %w", err)
	}
	log.Info("vault store initialized", zap.String("path", path))
	return &IsolatedVault{BaseDir: path, Key: key}, nil
}

type VaultStore interface {
	LoadConfig(key string) (*internal_environment.EnvConfig, error)
	SaveConfig(key string, cfg *internal_environment.EnvConfig) error
	LoadGoldenHash(machine string) (string, error)
	SealGoldenHash(machine string, hash []byte) error
	LoadFirstBootMarker() (*internal_boot.FirstBootMarker, error)
	MarkFirstBoot(*internal_boot.FirstBootMarker) error
	Read(collection, key string, out interface{}) (bool, error)
	Write(collection, key string, value interface{}) error
	Exists(collection, key string) (bool, error)
}

// Internal Storage Helpers (Encryption at Rest + Atomic Writes)

func (v *IsolatedVault) resolvePath(filename string) (string, error) {
	cleanPath := filepath.Clean(filepath.Join(v.BaseDir, filename))
	baseClean := filepath.Clean(v.BaseDir)

	if !strings.HasPrefix(cleanPath, baseClean+string(filepath.Separator)) && cleanPath != baseClean {
		return "", ErrPathTraversal
	}
	return cleanPath, nil
}

func (v *IsolatedVault) writeEncrypted(filename string, plaintext []byte) error {
	targetPath, err := v.resolvePath(filename)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0700); err != nil {
		return fmt.Errorf("failed to prepare vault directory: %w", err)
	}

	block, err := aes.NewCipher(v.Key)
	if err != nil {
		return fmt.Errorf("aes cipher initialization failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm mode initialization failed: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation failed: %w", err)
	}

	// Payload layout: [Nonce (12 Bytes)][Ciphertext + Authentication Tag]
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Atomic Write Strategy: write to temp file then rename
	tmpFile, err := os.CreateTemp(filepath.Dir(targetPath), ".vault-write-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tmpName := tmpFile.Name()

	if _, err := tmpFile.Write(ciphertext); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("failed writing encrypted data to temp file: %w", err)
	}
	_ = tmpFile.Close()

	if err := os.Chmod(tmpName, 0600); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("failed setting temp file permissions: %w", err)
	}

	if err := os.Rename(tmpName, targetPath); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("failed persisting record via rename: %w", err)
	}

	return nil
}

func (v *IsolatedVault) readDecrypted(filename string) ([]byte, error) {
	targetPath, err := v.resolvePath(filename)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	block, err := aes.NewCipher(v.Key)
	if err != nil {
		return nil, fmt.Errorf("aes cipher initialization failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm mode initialization failed: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return nil, errors.New("vault record is corrupted or improperly formatted")
	}

	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault payload: %w", err)
	}

	return plaintext, nil
}
