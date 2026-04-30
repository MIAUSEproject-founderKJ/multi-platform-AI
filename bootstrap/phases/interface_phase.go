//bootstraps/phases/interface_phase.go

package bootstrap_phase

import (
	"errors"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	audio_engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/audio/engine"
)

func PhaseInterface(ctx bootstrap.BootContext, caps *internal_environment.CapabilityProfile) (*user_setting.UserSession, error) {

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

func ToDeviceCapabilities(cp *internal_environment.CapabilityProfile) DeviceCapabilities {
	return DeviceCapabilities{
		HasDisplay:  cp.Has(internal_environment.CapDisplay),
		HasKeyboard: cp.Has(internal_environment.CapKeyboard),
		HasMic:      cp.Has(internal_environment.CapMicrophone),
		HasSpeaker:  cp.Has(internal_environment.CapSpeaker),
		GPU:         cp.Has(internal_environment.CapGPU),
	}
}

type MainInterface interface {
	Start(session *user_setting.UserSession) error
}

type HybridAuthUI struct {
	Voice audio_engine.VoiceEngine
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
