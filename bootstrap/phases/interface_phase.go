//bootstraps/phases/interface_phase.go

package bootstrap_phase

import (
	"errors"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

func PhaseInterface(ctx internal_boot.BootContext, caps *internal_verification.CapabilityProfile) (*user_setting.UserSession, error) {

	mode := ResolveInteractionMode(nil, caps.Set)

	ui := BuildAuthInterface(mode)
	if ui == nil {
		return nil, errors.New("failed to build auth interface")
	}

	authManager := auth.AuthManager{Vault: ctx.Vault}

	result, err := ui.StartAuthFlow(&authManager)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ToDeviceCapabilities(cp *internal_verification.CapabilityProfile) DeviceCapabilities {
	return DeviceCapabilities{
		HasDisplay:  cp.Has(internal_verification.CapDisplay),
		HasKeyboard: cp.Has(internal_verification.CapKeyboard),
		HasMic:      cp.Has(internal_verification.CapMicrophone),
		HasSpeaker:  cp.Has(internal_verification.CapSpeaker),
		GPU:         cp.Has(internal_verification.CapGPU),
	}
}

type MainInterface interface {
	Start(session *user_setting.UserSession) error
}

type HybridAuthUI struct {
	Voice runtime.VoiceEngine
	GUI   runtime.GUIEngine
}

func (h *HybridAuthUI) StartAuthFlow(auth *auth.AuthManager) (*user_setting.UserSession, error) {

	choice := h.promptChoice()

	switch choice {
	case "login":
		creds := h.collectCredentials()
		return auth.Login(creds.UserID, creds.Password)

	case "signup":
		return auth.Register()

	default:
		return nil, errors.New("invalid choice")
	}
}

type Credentials struct {
	UserID   string
	Password string
}

func (h *HybridAuthUI) promptChoice() string {
	fmt.Println("Choose: login / signup")
	var input string
	fmt.Scanln(&input)
	return input
}

func (h *HybridAuthUI) collectCredentials() Credentials {
	var user, pass string
	fmt.Print("User: ")
	fmt.Scanln(&user)
	fmt.Print("Pass: ")
	fmt.Scanln(&pass)

	return Credentials{
		UserID:   user,
		Password: pass,
	}
}
