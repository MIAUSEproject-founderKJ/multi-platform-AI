// mutual_interaction/interaction_mode_resolver.go
package mutual_interaction

import (
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

func ResolveInteractionMode(
	cfg *user_setting.CustomizedConfig,
	cap internal_environment.CapabilitySet,
) user_setting.InteractionMode {

	if cfg != nil && cfg.PreferredMode != "" && cfg.PreferredMode != "auto" {
		return user_setting.InteractionMode(cfg.PreferredMode)
	}

	switch {
	case cap.Has(internal_environment.CapDisplay) && cap.Has(internal_environment.CapGPU):
		return user_setting.ModeGUI

	case cap.Has(internal_environment.CapDisplay) && cap.Has(internal_environment.CapKeyboard):
		return user_setting.ModeTUI

	default:
		return user_setting.ModeCLI
	}
}
