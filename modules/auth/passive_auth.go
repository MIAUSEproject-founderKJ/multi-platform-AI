//modules/auth/passive_auth.go

package auth

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type PassiveAuthManager struct {
	Kernel *core.Kernel
}

// ResolveUser waits for a physical "Presence" signal instead of a password
func (pam *PassiveAuthManager) ResolveUser() (*Session, error) {
	logging.Info("[AUTH] Waiting for Physical Credential (NFC/Biometric)...")

	for {
		// 1. Check for NFC Badge
		if badgeID, found := pam.Kernel.Hardware.ScanNFC(); found {
			return pam.LoginByToken(badgeID, "NFC_TOKEN")
		}

		// 2. Check for Face ID (Biometric)
		if faceHash, found := pam.Kernel.Vision.DetectAuthorizedFace(); found {
			return pam.LoginByToken(faceHash, "BIOMETRIC")
		}

		// Prevent CPU pegging
		time.Sleep(500 * time.Millisecond)
	}
}