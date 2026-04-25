//boot/phases/interface_phase.go

package boot_phase

import (
	"errors"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/interaction"
	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
)

func PhaseInterface(ctx schema_boot.BootContext, caps *schema_security.CapabilityProfile) (*schema_identity.UserSession, error) {

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

func ToDeviceCapabilities(cp *schema_security.CapabilityProfile) DeviceCapabilities {
	return DeviceCapabilities{
		HasDisplay:  cp.Has(schema_security.CapDisplay),
		HasKeyboard: cp.Has(schema_security.CapKeyboard),
		HasMic:      cp.Has(schema_security.CapMicrophone),
		HasSpeaker:  cp.Has(schema_security.CapSpeaker),
		GPU:         cp.Has(schema_security.CapGPU),
	}
}

type MainInterface interface {
	Start(session *schema_identity.UserSession) error
}

type HybridAuthUI struct {
	Voice runtime.VoiceEngine
	GUI   runtime.GUIEngine
}

func (h *HybridAuthUI) StartAuthFlow(auth *auth.AuthManager) (*schema_identity.UserSession, error) {

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
