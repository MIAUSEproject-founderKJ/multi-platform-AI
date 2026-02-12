//core/kernel_commands.go
package core

import (
    "fmt" // Added for Errorf
    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/commands"
    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

// ProcessCommand decomposes a high-level task into hardware actions.
func (k *Kernel) ProcessCommand(cmd commands.Task) error {
    logging.Info("[COMMAND] Received: %s (Priority: %d)", cmd.Type, cmd.Priority)

    // 1. SECURITY GATE: Verify we are in the correct mode
    if k.Trust.OperationMode == "MANUAL_ONLY" && cmd.Type != commands.CmdHalt {
        logging.Error("[COMMAND] Blocked: System in Manual-Only mode due to low trust (Score: %.2f)", k.Trust.CurrentScore)
        return fmt.Errorf("command_forbidden: low_trust_state")
    }

    // 2. ROUTING
    switch cmd.Type {
    case commands.CmdHalt:
        // Immediate hardware stop
        k.Bridge.TransitionTo("EMERGENCY_OFF")
        return nil

    case commands.CmdNavigate:
        return k.handleNavigation(cmd.Params)

    case commands.CmdScan:
        return k.handlePerceptionScan(cmd.Params)

    default:
        return fmt.Errorf("unknown_command_type: %s", cmd.Type)
    }
}

func (k *Kernel) handleNavigation(params map[string]interface{}) error {
    lat, okLat := params["lat"].(float64)
    lng, okLng := params["lng"].(float64)
    
    if !okLat || !okLng {
        return fmt.Errorf("invalid_navigation_params: lat/lng missing or malformed")
    }

    logging.Info("[NAV] Plotting course to %f, %f", lat, lng)
    
    // Command the bridge to move actuators
    // Note: PowerController interface must now include WriteActuator
    return k.Bridge.WriteActuator("drive_motor", 0.5) 
}

func (k *Kernel) handlePerceptionScan(params map[string]interface{}) error {
    sensorType, _ := params["sensor"].(string)
    logging.Info("[SCAN] Initiating perception sweep using %s", sensorType)
    
    // Trigger a vision frame capture and processing
    frame := k.Hardware.GetCameraFrame()
    results := k.Vision.ProcessFrame(frame)
    
    // Store result in Cognitive Vault for "Dream State" recall later
    k.Memory.Store("last_perception_scan", results)
    
    return nil
}