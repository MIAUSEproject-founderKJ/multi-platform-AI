//runtime/mode_resolver.go

package runtime

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

type InteractionMode string

const (
	ModeCLI   InteractionMode = "cli"
	ModeTUI   InteractionMode = "tui"
	ModeGUI   InteractionMode = "gui"
	ModeVoice InteractionMode = "voice"
)

func ResolveInteractionMode(
	cfg *schema.CustomizedConfig,
	cap schema.CapabilitySet,
) InteractionMode {

	if cfg != nil && cfg.PreferredMode != "" && cfg.PreferredMode != "auto" {
		return InteractionMode(cfg.PreferredMode)
	}

	switch {
	case cap.Has(schema.CapDisplay) && cap.Has(schema.CapGPU):
		return ModeGUI

	case cap.Has(schema.CapDisplay) && cap.Has(schema.CapKeyboard):
		return ModeTUI

	case cap.Has(schema.CapMicrophone) && cap.Has(schema.CapSpeaker):
		return ModeVoice

	default:
		return ModeCLI
	}
}