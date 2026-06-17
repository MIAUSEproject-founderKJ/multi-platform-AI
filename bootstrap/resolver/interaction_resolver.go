// bootstrap/resolver/interaction_resolver.go
package bootstrap_resolver

import (
	"context"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type InterfaceAdapter interface {
	Start(session *user_setting.UserSession) error
	Stop(ctx context.Context) error
	Notify(msg string)
}

func ResolveInteractionCapability(
	dev internal_environment.DeviceCapabilities,
) user_setting.InteractionCapability {

	return user_setting.InteractionCapability{
		CLI: true,

		TUI: dev.HasDisplay &&
			dev.HasKeyboard,

		GUI: dev.HasDisplay &&
			dev.GPU,

		Voice: dev.HasMic &&
			dev.HasSpeaker,
	}
}

func ResolvePreferredInteractionMode(
	cap user_setting.InteractionCapability,
) user_setting.InteractionMode {

	switch {

	case cap.GUI && cap.TUI && cap.Voice:
		return user_setting.ModeFull

	case cap.GUI && cap.Voice:
		return user_setting.ModeGTVio

	case cap.GUI:
		return user_setting.ModeGUIonly

	case cap.TUI:
		return user_setting.ModeTUIonly

	case cap.Voice:
		return user_setting.ModeVio

	default:
		return user_setting.ModeCLIonly
	}
}
