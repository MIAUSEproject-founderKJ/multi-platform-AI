//boot/phases/capability_phase.go

package boot_phase

import (
	"fmt"
	"os"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
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

func BuildAuthInterface(mode schema_identity.InteractionMode) auth.AuthInterface {
	switch mode {
	case schema_identity.ModeGUI:
		return NewGUIAuth()
	case schema_identity.ModeTUI:
		return NewTUIAuth()
	case schema_identity.ModeVoice:
		return NewVoiceAuth()
	default:
		return NewCLIAuth()
	}
}

func ResolveInteractionMode(
	cfg *schema_security.CustomizedConfig,
	cap schema_security.CapabilitySet,
) schema_identity.InteractionMode {

	if cfg != nil && cfg.PreferredMode != "" && cfg.PreferredMode != "auto" {
		return schema_identity.InteractionMode(cfg.PreferredMode)
	}

	switch {
	case cap.Has(schema_security.CapDisplay) && cap.Has(schema_security.CapGPU):
		return schema_identity.ModeGUI

	case cap.Has(schema_security.CapDisplay) && cap.Has(schema_security.CapKeyboard):
		return schema_identity.ModeTUI

	case cap.Has(schema_security.CapMicrophone) && cap.Has(schema_security.CapSpeaker):
		return schema_identity.ModeVoice

	default:
		return schema_identity.ModeCLI
	}
}

func PhaseCapability() *schema_security.CapabilityProfile {
	cp := schema_security.NewCapabilityProfile()

	// ---- Display ----
	if !isHeadless() {
		cp.Mark(schema_security.CapDisplay, schema_security.CapOK)
	} else {
		cp.Mark(schema_security.CapDisplay, schema_security.CapUnavailable)
	}

	// ---- Keyboard ----
	cp.Mark(schema_security.CapKeyboard, schema_security.CapOK) // assume present

	// ---- Network ----
	if hasNetwork() {
		cp.Mark(schema_security.CapNetwork, schema_security.CapOK)
	} else {
		cp.Mark(schema_security.CapNetwork, schema_security.CapDegraded)
	}

	// ---- Microphone (stub) ----
	if hasMicrophone() {
		cp.Mark(schema_security.CapMicrophone, schema_security.CapOK)
	} else {
		cp.Mark(schema_security.CapMicrophone, schema_security.CapUnavailable)
	}

	// ---- Speaker (stub) ----
	if hasSpeaker() {
		cp.Mark(schema_security.CapSpeaker, schema_security.CapOK)
	} else {
		cp.Mark(schema_security.CapSpeaker, schema_security.CapUnavailable)
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
	Start(session *schema_identity.UserSession) error
	Notify(msg string)
}

type DeviceCapabilities struct {
	HasDisplay  bool
	HasKeyboard bool
	HasMic      bool
	HasSpeaker  bool
	GPU         bool
}

func SelectInteractionMode(cap DeviceCapabilities) schema_identity.InteractionMode {

	switch {
	case cap.HasDisplay && cap.GPU:
		return schema_identity.ModeGUI

	case cap.HasDisplay && cap.HasKeyboard:
		return schema_identity.ModeTUI

	case cap.HasMic && cap.HasSpeaker:
		return schema_identity.ModeVoice

	default:
		return schema_identity.ModeCLI
	}
}

// =========================================
type CLIAdapter struct{}

func (c *CLIAdapter) Start(session *schema_identity.UserSession) error {
	fmt.Println("CLI session started:", session.SessionID)
	return nil
}

// =========================================
type TUIAdapter struct{}

func (t *TUIAdapter) Start(session *schema_identity.UserSession) error {
	// integrate charmbracelet/bubbletea
	return nil
}

// =========================================
type GUIAdapter struct{}

func (g *GUIAdapter) Start(session *schema_identity.UserSession) error {
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

func (o *Orchestrator) StartAll(s *schema_identity.UserSession) {
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

func (s *ScreenAdapter) Start(session *schema_identity.UserSession) error {
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

func (c *CLIAuth) StartAuthFlow(auth *auth.AuthManager) (*schema_identity.UserSession, error) {
	return auth.LoginOrSignUpInteractive()
}
func (t *TUIAuth) StartAuthFlow(am *auth.AuthManager) (*schema_identity.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}

func (g *GUIAuth) StartAuthFlow(am *auth.AuthManager) (*schema_identity.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}

func (v *VoiceAuth) StartAuthFlow(am *auth.AuthManager) (*schema_identity.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}
