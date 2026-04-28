//core/verification/identity/token_service.go

package verification_identity

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
