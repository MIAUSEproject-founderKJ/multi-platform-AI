//core/security/persistence/key_manager.go

package verification_persistence

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
)

const RequiredKeyLength = 32 // 256-bit key for AES-256

// GenerateSecureKeyBase64 creates a cryptographically secure 32-byte key encoded in Base64.
func GenerateSecureKeyBase64() (string, error) {
	key := make([]byte, RequiredKeyLength)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// LoadSecureKey fetches and validates the encryption key from environment variables.
// If missing, generates an ephemeral key for temporary execution.
func LoadSecureKey() ([]byte, error) {
	keyStr := os.Getenv("APP_ENCRYPTION_KEY")

	if keyStr == "" {
		logging.Warn("[VAULT] APP_ENCRYPTION_KEY not set; generating ephemeral key. Data will not persist across restarts.")

		gen, err := GenerateSecureKeyBase64()
		if err != nil {
			return nil, fmt.Errorf("ephemeral key generation failed: %w", err)
		}

		keyBytes, _ := base64.StdEncoding.DecodeString(gen)
		return keyBytes, nil
	}

	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 encryption key: %w", err)
	}

	if len(key) != RequiredKeyLength {
		return nil, fmt.Errorf("invalid key length: got %d bytes, expected %d bytes", len(key), RequiredKeyLength)
	}

	return key, nil
}
