//boot/phase_auth_interface.go

package boot

func PhaseAuthInterface(ctx BootContext, caps *schema.CapabilityProfile) (*schema.UserSession, error) {

	mode := schema.SelectInteractionMode(caps)

	ui := schema.BuildAuthInterface(mode)
	if ui == nil {
		return nil, errors.New("failed to build auth interface")
	}

	authManager := AuthManager{Vault: ctx.Vault}

	result, err := ui.StartAuthFlow(authManager)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type AuthInterface interface {
	StartAuthFlow(auth AuthManager) (*schema.UserSession, error)
}

type MainInterface interface {
	Start(session *schema.UserSession) error
}

type HybridAuthUI struct {
	Voice VoiceEngine
	GUI   GUIEngine
}

func (h *HybridAuthUI) StartAuthFlow(auth AuthManager) (*schema.UserSession, error) {

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