//core/kernel_commands.go

package core

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/commands"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/bridge"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

// ProcessCommand decomposes a high-level task into hardware actions.
func (k *Kernel) ProcessCommand(cmd commands.Task) error {
	logging.Info("[COMMAND] Received: %s (Priority: %d)", cmd.Type, cmd.Priority)

	// 1. SECURITY GATE: Verify we are in the correct mode
	if k.Trust.OperationMode == "MANUAL_ONLY" && cmd.Type != commands.CmdHalt {
		logging.Error("[COMMAND] Blocked: System in Manual-Only mode due to low trust.")
		return fmt.Errorf("command_forbidden: low_trust_state")
	}

	// 2. ROUTING: Send to the appropriate subsystem
	switch cmd.Type {
	case commands.CmdHalt:
		k.Bridge.TransitionTo(bridge.StateEmergencyOff)
		return nil

	case commands.CmdNavigate:
		return k.handleNavigation(cmd.Params)

	case commands.CmdScan:
		return k.handlePerceptionScan(cmd.Params)
	}

	return nil
}

func (k *Kernel) handleNavigation(params map[string]interface{}) error {
	// Example: Extract coordinates and send to the Bridge
	lat := params["lat"]
	lng := params["lng"]
	logging.Info("[NAV] Plotting course to %v, %v", lat, lng)
	
	// Bridge would then be used to move the physical motors
	return k.Bridge.WriteActuator("drive_motor", 0.5)
}