//interaction/interaction_mode.go

package interaction

import (
	"fmt"
	"os"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
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

func DetectCapabilityProfile() *schema.CapabilityProfile {
	cp := schema.NewCapabilityProfile()

	// ---- Display ----
	if !isHeadless() {
		cp.Mark(schema.CapDisplay, schema.CapOK)
	} else {
		cp.Mark(schema.CapDisplay, schema.CapUnavailable)
	}

	// ---- Keyboard ----
	cp.Mark(schema.CapKeyboard, schema.CapOK) // assume present

	// ---- Network ----
	if hasNetwork() {
		cp.Mark(schema.CapNetwork, schema.CapOK)
	} else {
		cp.Mark(schema.CapNetwork, schema.CapDegraded)
	}

	// ---- Microphone (stub) ----
	if hasMicrophone() {
		cp.Mark(schema.CapMicrophone, schema.CapOK)
	} else {
		cp.Mark(schema.CapMicrophone, schema.CapUnavailable)
	}

	// ---- Speaker (stub) ----
	if hasSpeaker() {
		cp.Mark(schema.CapSpeaker, schema.CapOK)
	} else {
		cp.Mark(schema.CapSpeaker, schema.CapUnavailable)
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
	Start(session *schema.UserSession) error
	Notify(msg string)
}

type DeviceCapabilities struct {
	HasDisplay  bool
	HasKeyboard bool
	HasMic      bool
	HasSpeaker  bool
	GPU         bool
}
type InteractionMode string

const (
	ModeCLI   InteractionMode = "cli"
	ModeTUI   InteractionMode = "tui"
	ModeGUI   InteractionMode = "gui"
	ModeVoice InteractionMode = "voice"
)

func SelectInteractionMode(cap DeviceCapabilities) InteractionMode {

	switch {
	case cap.HasDisplay && cap.GPU:
		return ModeGUI

	case cap.HasDisplay && cap.HasKeyboard:
		return ModeTUI

	case cap.HasMic && cap.HasSpeaker:
		return ModeVoice

	default:
		return ModeCLI
	}
}

// =========================================
type CLIAdapter struct{}

func (c *CLIAdapter) Start(session *schema.UserSession) error {
	fmt.Println("CLI session started:", session.SessionID)
	return nil
}

// =========================================
type TUIAdapter struct{}

func (t *TUIAdapter) Start(session *schema.UserSession) error {
	// integrate charmbracelet/bubbletea
	return nil
}

// =========================================
type GUIAdapter struct{}

func (g *GUIAdapter) Start(session *schema.UserSession) error {
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
type WhisperSTT struct{}

func (w *WhisperSTT) Listen() (string, error) {
	return "parsed speech command", nil
}

type SystemTTS struct{}

func (s *SystemTTS) Speak(text string) error {
	fmt.Println("[TTS]:", text)
	return nil
}

// =========================================Adapter Factory (Core Integration)
func BuildInterface(mode InteractionMode) InterfaceAdapter {
	switch mode {
	case ModeGUI:
		return &GUIAdapter{}
	case ModeTUI:
		return &TUIAdapter{}
	case ModeVoice:
		return NewVoiceAdapter()
	default:
		return &CLIAdapter{}
	}
}

type Orchestrator struct {
	adapters []InterfaceAdapter
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Add(adapter InterfaceAdapter) {
	o.adapters = append(o.adapters, adapter)
}

func (o *Orchestrator) StartAll(session *schema.UserSession) {
	for _, a := range o.adapters {
		go a.Start(session)
	}
}

func (o *Orchestrator) Broadcast(msg string) {
	for _, a := range o.adapters {
		a.Notify(msg)
	}
}

func BuildOrchestrator(cp *schema.CapabilityProfile) *Orchestrator {
	orch := NewOrchestrator()

	if cp.IsHealthy(schema.CapDisplay) {
		orch.Add(NewScreenAdapter())
	}

	if cp.IsHealthy(schema.CapMicrophone) &&
		cp.IsHealthy(schema.CapSpeaker) {
		orch.Add(NewVoiceAdapter())
	}

	return orch
}

// ===================Screen Adapter
type ScreenAdapter struct{}

func NewScreenAdapter() *ScreenAdapter {
	return &ScreenAdapter{}
}

// ============VoiceAdapter
type VoiceAdapter struct {
	engine *VoiceEngine
}

func NewVoiceAdapter() *VoiceAdapter {
	return &VoiceAdapter{
	engine: &VoiceEngine{
		STT: &WhisperSTT{},
		TTS: &SystemTTS{},
	},
}
}

func (s *ScreenAdapter) Start(session *schema.UserSession) error {
	fmt.Println("Screen adapter started")
	return nil
}

func (s *ScreenAdapter) Notify(msg string) {
	fmt.Println("[SCREEN]", msg)
}

func (v *VoiceAdapter) Notify(msg string) {
	v.engine.outputChan <- msg
}

func (v *VoiceAdapter) Start(session *schema.UserSession) error {
	v.engine.Start()
	return nil
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

type VoiceEngine struct {
	STT        SpeechToText
	TTS        TextToSpeech
	outputChan chan string
}

func NewVoiceEngine() *VoiceEngine {
	return &VoiceEngine{
		STT:        &WhisperSTT{},
		TTS:        &SystemTTS{},
		outputChan: make(chan string, 10),
	}
}

func (v *VoiceEngine) Start() {
	go func() {
		for msg := range v.outputChan {
			v.TTS.Speak(msg)
		}
	}()
}
