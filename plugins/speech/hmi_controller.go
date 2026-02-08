//MIAUSEproject-founderKJ/multi-platform-AI/plugins/speech/hmi_controller.go (human_user_interface controller)

package speech

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type HMIController struct {
	Config *defaults.EnvConfig
	Mode   string // "Visual", "Auditory", "Headless"
}

func NewHMIController(env *defaults.EnvConfig) *HMIController {
	hmi := &HMIController{Config: env}
	hmi.determineModality()
	return hmi
}

func (h *HMIController) determineModality() {
	// 1. REFLECTIVE: Check hardware for display capabilities
	hasScreen := false
	for _, bus := range h.Config.Hardware.Buses {
		if bus.Type == "hdmi" || bus.Type == "edp" || bus.Type == "usb_c_display" {
			hasScreen = true
			break
		}
	}

	// 2. RESPONSIVE: Set mode
	if hasScreen && h.Config.Platform.Final != defaults.PlatformEmbedded {
		h.Mode = "Visual"
	} else if h.Config.Hardware.Battery.Present {
		h.Mode = "Auditory" // Save power, use voice only
	} else {
		h.Mode = "Headless"
	}

	logging.Info("[HMI] Selected Modality: %s", h.Mode)
}
