//core\security\persistence\key_manager.go

package security_persistence

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
)

// KeyManager handles the generation, storage, and retrieval of cryptographic keys.
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
