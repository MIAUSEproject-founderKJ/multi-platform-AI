// modules\interfaces\CLIAuth.go
package interfaces

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"

func BuildAuthInterface(mode user_setting.InteractionMode) auth.AuthInterface {
	switch mode {
	case user_setting.ModeGUI:
		return gui.NewAuth()
	case user_setting.ModeTUI:
		return tui.NewAuth()
	case user_setting.ModeVoice:
		return voice.NewAuth()
	default:
		return cli.NewAuth()
	}
}
