//bridge/power_controller.go
package bridge

import (
	"fmt"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/hmi"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/policy"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type PowerState string

const (
	StateFullPower   PowerState = "ACTUATORS_ON"
	StateLowPower    PowerState = "SENSORS_ONLY"
	StateEmergencyOff PowerState = "ISOLATED"
)

type BridgeController struct {
	CurrentState PowerState
	ActiveBuses  []string
}

// ExecuteCommand acts as the "Gatekeeper" for physical movement.
func (bc *BridgeController) SyncWithTrust(trust *policy.TrustDescriptor) {
	logging.Info("[BRIDGE] Syncing Physical State with Trust: %.2f", trust.CurrentScore)

	switch trust.OperationMode {
	case "AUTONOMOUS":
		bc.TransitionTo(StateFullPower)
	case "ASSISTED":
		bc.TransitionTo(StateLowPower)
	case "MANUAL_ONLY":
		bc.TransitionTo(StateEmergencyOff)
	}
}

func (bc *BridgeController) TransitionTo(target PowerState) {
	if bc.CurrentState == target {
		return
	}
	
	logging.Warn("[BRIDGE] Hardware Transition: %s -> %s", bc.CurrentState, target)
	bc.CurrentState = target
	
	// Real-world logic here would involve writing to GPIO pins or 
	// sending "Disable" frames to Motor Controllers via CAN-bus.
}