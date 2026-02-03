//MIAUSEproject-founderKJ/multi-platform-AI/bridge/control/actuator_loop.go

//This file handles the high-frequency translation of AI "Intents" into hardware "Commands." It lives in the bridge/ layer because it spans the gap between the logic of the core and the wires of the hal.

package control

import (
	"context"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bridge/hal"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/scheduler"
)

type ActuatorLoop struct {
	interlock *hal.SafetyInterlock
	bus       *hal.BusProvider // Interface for CAN, USB, or GPIO
	hz        int              // Target frequency (e.g., 100Hz)
}

// Start initiates the deterministic control cycle
func (al *ActuatorLoop) Start(ctx context.Context) {
	// Use the internal scheduler for real-time priority
	ticker := scheduler.NewPreciseTicker(al.hz)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			al.tick()
		}
	}
}

func (al *ActuatorLoop) tick() {
	// 1. GATING: Check the Safety Interlock status BEFORE processing
	if atomic.LoadInt32(&al.interlock.CurrentState) != int32(hal.StateClear) {
		// Interlock is active; force Neutral/Safe state and skip AI command
		al.bus.WriteSafeState()
		return
	}

	// 2. FETCH: Get the next movement intent from the Cognition/Navigation layer
	// This is a non-blocking read from a ring buffer
	intent := al.fetchLatestIntent()

	// 3. TRANSLATE: Convert floating point AI intent to Raw Bus Values (e.g. 0-255)
	rawCommand := al.translateIntentToHardware(intent)

	// 4. EXECUTE: Physical write to the hardware bus
	al.bus.Write(rawCommand)
}
