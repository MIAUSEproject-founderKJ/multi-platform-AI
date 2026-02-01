//MIAUSEproject-founderKJ/multi-platform-AI/plugins/speech/hmi_controller.go (human_user_interface controller)

package speech

import (
	"multi-platform-AI/configs/configStruct"
	"multi-platform-AI/internal/logging"
)

type HMIController struct {
	Config *configStruct.EnvConfig
	Mode   string // "Visual", "Auditory", "Headless"
}

func NewHMIController(env *configStruct.EnvConfig) *HMIController {
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
	if hasScreen && h.Config.Platform.Final != configStruct.PlatformEmbedded {
		h.Mode = "Visual"
	} else if h.Config.Hardware.Battery.Present {
		h.Mode = "Auditory" // Save power, use voice only
	} else {
		h.Mode = "Headless"
	}
	
	logging.Info("[HMI] Selected Modality: %s", h.Mode)
}