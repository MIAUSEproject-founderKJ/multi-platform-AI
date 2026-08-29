// modules/authui/hybrid.go
package authui

import (
	"errors"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type HybridAuthUI struct {
	Voice VoiceEngine
	GUI   GUIEngine
}

func (h *HybridAuthUI) StartAuthFlow(
	manager *auth.AuthManager,
) (*user_setting.UserSession, error) {

	choice := h.promptChoice()

	switch choice {
	case "login":
		creds := h.collectCredentials()
		return manager.Login(creds.UserID, creds.Password)

	case "signup":
		return manager.Register()

	default:
		return nil, errors.New("invalid choice")
	}
}
