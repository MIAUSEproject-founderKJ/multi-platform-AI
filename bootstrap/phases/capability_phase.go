//bootstrap/phases/capability_phase.go

package bootstrap_phase

import (
	"context"
	"os"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

func PhaseCapability() *internal_environment.CapabilityProfile {
	cp := internal_environment.NewCapabilityProfile()

	// ---- Display ----
	if !isHeadless() {
		cp.Mark(internal_environment.CapDisplay, internal_environment.CapOK)
	} else {
		cp.Mark(internal_environment.CapDisplay, internal_environment.CapUnavailable)
	}

	// ---- Keyboard ----
	cp.Mark(internal_environment.CapKeyboard, internal_environment.CapOK) // assume present

	// ---- Network ----
	if hasNetwork() {
		cp.Mark(internal_environment.CapNetwork, internal_environment.CapOK)
	} else {
		cp.Mark(internal_environment.CapNetwork, internal_environment.CapDegraded)
	}

	// ---- Microphone (stub) ----
	if hasMicrophone() {
		cp.Mark(internal_environment.CapMicrophone, internal_environment.CapOK)
	} else {
		cp.Mark(internal_environment.CapMicrophone, internal_environment.CapUnavailable)
	}

	// ---- Speaker (stub) ----
	if hasSpeaker() {
		cp.Mark(internal_environment.CapSpeaker, internal_environment.CapOK)
	} else {
		cp.Mark(internal_environment.CapSpeaker, internal_environment.CapUnavailable)
	}

	return cp
}

func isHeadless() bool {
	return os.Getenv("DISPLAY") == "" &&
		os.Getenv("WAYLAND_DISPLAY") == ""
}

func hasNetwork() bool {
	// simple heuristic
	return true
}

func hasMicrophone() bool {
	// stub → replace with OS API later
	return true
}

func hasSpeaker() bool {
	return true
}

type InterfaceAdapter interface {
	Start(session *user_setting.UserSession) error
	Stop(ctx context.Context) error
	Notify(msg string)
}

func SelectInteractionMode(cap internal_environment.DeviceCapabilities) user_setting.InteractionMode {

	switch {
	case cap.HasDisplay && cap.GPU:
		return user_setting.ModeGUI

	case cap.HasDisplay && cap.HasKeyboard:
		return user_setting.ModeTUI

	case cap.HasMic && cap.HasSpeaker:
		return user_setting.ModeVoice

	default:
		return user_setting.ModeCLI
	}
}

type SpeechToText interface {
	Listen() (string, error)
}

type TextToSpeech interface {
	Speak(string) error
}
