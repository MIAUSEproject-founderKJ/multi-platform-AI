// bridge/hal/types.go
package hal

import (
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

// InterlockState uses iota for clean incrementing constants
type InterlockState int32

const (
	StateClear       InterlockState = 0
	StateEStop       InterlockState = 1
	StateWarning     InterlockState = 2
	StateObstruction InterlockState = 3
)

// RawCommand represents a single packet destined for the hardware bus
type RawCommand struct {
	ID   uint32
	Data []byte
}

// BusProvider interface for hardware abstraction
type BusProvider interface {
	Write(cmd RawCommand) error
	WriteSafeState() error
}

// SafetyInterlock handles hardware-level "Kill Switches"
type SafetyInterlock struct {
	// Must be int32 for atomic operations.
	// We handle the "string" representation via the GetStatus() method.
	currentState int32
}

// PollHardware is called by the Real-Time Scheduler
func (si *SafetyInterlock) PollHardware() {
	rawStatus := si.checkPhysicalEStopPin()

	if rawStatus == StateEStop {
		si.TriggerEmergencyStop("Physical E-Stop Pressed")
	}
}

func (si *SafetyInterlock) TriggerEmergencyStop(reason string) {
	// 1. Atomically update state
	atomic.StoreInt32(&si.currentState, int32(StateEStop))

	// 2. Log with the required signature
	slog.Error("CRITICAL_SAFETY_INTERVENTION", "reason", reason)

	// Fix: Passing an actual error to the Protect function as the compiler requested
	errReason := fmt.Errorf("safety intervention: %s", reason)
	logging.ProtectIncidentLog(errReason)

	// 3. Hardware Kill
	si.killAllActuators()
}

func (si *SafetyInterlock) killAllActuators() {
	// Implementation for Neutral-Gear / 0-Voltage signals
}

// GetStatus returns the human-readable string for the HMI/UI
func (si *SafetyInterlock) GetStatus() string {
	state := InterlockState(atomic.LoadInt32(&si.currentState))
	switch state {
	case StateClear:
		return "CLEAR"
	case StateEStop:
		return "EMERGENCY_STOP"
	case StateWarning:
		return "WARNING"
	case StateObstruction:
		return "OBSTRUCTION"
	default:
		return "UNKNOWN"
	}
}

// Private helper to simulate hardware pin reading
func (si *SafetyInterlock) checkPhysicalEStopPin() InterlockState {
	// Actual GPIO/CAN logic goes here
	return StateClear
}
