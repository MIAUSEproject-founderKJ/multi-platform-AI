//internal/schema/interaction_mode.go

package schema



type DeviceCapabilities struct {
	HasDisplay bool
	HasKeyboard bool
	HasMic bool
	HasSpeaker bool
	GPU bool
}

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


//=========================================
type CLIAdapter struct{}

func (c *CLIAdapter) Start(session *schema.UserSession) error {
	fmt.Println("CLI session started:", session.SessionID)
	return nil
}
//=========================================
type TUIAdapter struct{}

func (t *TUIAdapter) Start(session *schema.UserSession) error {
	// integrate charmbracelet/bubbletea
	return nil
}

//=========================================
type GUIAdapter struct{}

func (g *GUIAdapter) Start(session *schema.UserSession) error {
	// Launch window
	return nil
}
//=========================================
type VoiceEngine struct {
	STT SpeechToText
	TTS TextToSpeech
}

type SpeechToText interface {
	Listen() (string, error)
}

type TextToSpeech interface {
	Speak(string) error
}
//=========================================Example (Whisper + OS TTS)
type WhisperSTT struct{}

func (w *WhisperSTT) Listen() (string, error) {
	return "parsed speech command", nil
}

type SystemTTS struct{}

func (s *SystemTTS) Speak(text string) error {
	fmt.Println("[TTS]:", text)
	return nil
}

//=========================================Adapter Factory (Core Integration)
func BuildInterface(mode InteractionMode) InterfaceAdapter {
	switch mode {
	case ModeGUI:
		return &GUIAdapter{}
	case ModeTUI:
		return &TUIAdapter{}
	case ModeVoice:
		return &VoiceAdapter{
			Engine: &VoiceEngine{
				STT: &WhisperSTT{},
				TTS: &SystemTTS{},
			},
		}
	default:
		return &CLIAdapter{}
	}
}

//=========================================End-to-End Flow (Final Runtime)
boot := BootEngine{Vault: vault, Logger: logger}
seq, _ := boot.Initialize(env)

auth := AuthManager{Vault: vault}
session, _ := auth.Login()

cap := DetectDeviceCapabilities()
mode := SelectInteractionMode(cap)

ui := BuildInterface(mode)
ui.Start(session)