// bridge/actuators.go

func (bc *BridgeController) WriteActuator(busID string, signal float64) error {
	// SAFETY CHECK: Never move if the trust gate is closed
	if bc.CurrentState == StateEmergencyOff {
		return fmt.Errorf("safety_interlock_active: movement_forbidden")
	}

	// Logic to translate 0.0-1.0 signal to 16-bit PWM or CAN-bus frames
	logging.Debug("[BRIDGE] Writing to %s: Signal %.2f", busID, signal)
	return nil
}