//bootstrap/phases/capability_phase.go

package bootstrap_phase

import (
	"fmt"
	"os"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

type CLIAuth struct{}
type TUIAuth struct{}
type GUIAuth struct{}
type VoiceAuth struct{}

func NewCLIAuth() auth.AuthInterface {
	return &CLIAuth{}
}

func NewTUIAuth() auth.AuthInterface {
	return &TUIAuth{}
}

func NewGUIAuth() auth.AuthInterface {
	return &GUIAuth{}
}

func NewVoiceAuth() auth.AuthInterface {
	return &VoiceAuth{}
}

func (c *CLIAuth) Authenticate() error {
	return nil
}

func (t *TUIAuth) Authenticate() error {
	return nil
}

func (g *GUIAuth) Authenticate() error {
	return nil
}

func (v *VoiceAuth) Authenticate() error {
	return nil
}

func BuildAuthInterface(mode user_setting.InteractionMode) auth.AuthInterface {
	switch mode {
	case user_setting.ModeGUI:
		return NewGUIAuth()
	case user_setting.ModeTUI:
		return NewTUIAuth()
	case user_setting.ModeVoice:
		return NewVoiceAuth()
	default:
		return NewCLIAuth()
	}
}

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

	case cap.Has(internal_environment.CapMicrophone) && cap.Has(internal_environment.CapSpeaker):
		return user_setting.ModeVoice

	default:
		return user_setting.ModeCLI
	}
}

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
	Notify(msg string)
}

type DeviceCapabilities struct {
	HasDisplay  bool
	HasKeyboard bool
	HasMic      bool
	HasSpeaker  bool
	GPU         bool
}

func SelectInteractionMode(cap DeviceCapabilities) user_setting.InteractionMode {

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

// =========================================
type CLIAdapter struct{}

func (c *CLIAdapter) Start(session *user_setting.UserSession) error {
	fmt.Println("CLI session started:", user_setting.UserIdentity)
	return nil
}

// =========================================
type TUIAdapter struct{}

func (t *TUIAdapter) Start(session *user_setting.UserSession) error {
	// integrate charmbracelet/bubbletea
	return nil
}

// =========================================
type GUIAdapter struct{}

func (g *GUIAdapter) Start(session *user_setting.UserSession) error {
	// Launch window
	return nil
}

//=========================================

type SpeechToText interface {
	Listen() (string, error)
}

type TextToSpeech interface {
	Speak(string) error
}

// =========================================Example (Whisper + OS TTS)

// =========================================Adapter Factory (Core Integration)

type Orchestrator struct {
	adapters []InterfaceAdapter
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Add(adapter InterfaceAdapter) {
	o.adapters = append(o.adapters, adapter)
}

func (o *Orchestrator) StartAll(s *user_setting.UserSession) {
	for _, a := range o.adapters {
		go a.Start(s)
	}
}

func (o *Orchestrator) Broadcast(msg string) {
	for _, a := range o.adapters {
		a.Notify(msg)
	}
}

// ===================Screen Adapter
type ScreenAdapter struct{}

func NewScreenAdapter() *ScreenAdapter {
	return &ScreenAdapter{}
}

// ============VoiceAdapter

func (s *ScreenAdapter) Start(session *user_setting.UserSession) error {
	fmt.Println("Screen adapter started")
	return nil
}

func (s *ScreenAdapter) Notify(msg string) {
	fmt.Println("[SCREEN]", msg)
}

func (c *CLIAdapter) Notify(msg string) {
	fmt.Println("[CLI]", msg)
}

func (c *TUIAdapter) Notify(msg string) {
	fmt.Println("[TUI]", msg)
}

func (c *GUIAdapter) Notify(msg string) {
	fmt.Println("[GUI]", msg)
}

func (c *CLIAuth) StartAuthFlow(am *auth.AuthManager) (*user_setting.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}
func (t *TUIAuth) StartAuthFlow(am *auth.AuthManager) (*user_setting.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}

func (g *GUIAuth) StartAuthFlow(am *auth.AuthManager) (*user_setting.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}

func (v *VoiceAuth) StartAuthFlow(am *auth.AuthManager) (*user_setting.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}
