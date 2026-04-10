//boot/phase_auth_interface.go

package boot

import (
	"errors"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/interaction"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func ToDeviceCapabilities(cp *schema.CapabilityProfile) interaction.DeviceCapabilities {
	return interaction.DeviceCapabilities{
		HasDisplay:  cp.Has(schema.CapDisplay),
		HasKeyboard: cp.Has(schema.CapKeyboard),
		HasMic:      cp.Has(schema.CapMicrophone),
		HasSpeaker:  cp.Has(schema.CapSpeaker),
		GPU:         cp.Has(schema.CapGPU),
	}
}

func PhaseAuthInterface(ctx schema.BootContext, caps *schema.CapabilityProfile) (*schema.UserSession, error) {

	mode := interaction.SelectInteractionMode(ToDeviceCapabilities(caps))

	ui := schema.BuildAuthInterface(mode)
	if ui == nil {
		return nil, errors.New("failed to build auth interface")
	}

	authManager := auth.AuthManager{Vault: ctx.Vault}

	result, err := ui.StartAuthFlow(authManager)
	if err != nil {
		return nil, err
	}

	return result, nil
}



type MainInterface interface {
	Start(session *schema.UserSession) error
}

type HybridAuthUI struct {
	Voice interaction.VoiceEngine
	GUI   interaction.GUIEngine
}

func (h *HybridAuthUI) StartAuthFlow(auth auth.AuthManager) (*schema.UserSession, error) {

	choice := h.promptChoice()

	switch choice {
	case "login":
		creds := h.collectCredentials()
		return auth.Login(creds)

	case "signup":
		creds := h.collectCredentials()
		return auth.Register(creds)

	default:
		return nil, errors.New("invalid choice")
	}
}
