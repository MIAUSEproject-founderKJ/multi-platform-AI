//MIAUSEproject-founderKJ/multi-platform-AI/bridge/control/actuator_loop.go

//This file handles the high-frequency translation of AI "Intents" into hardware "Commands." It lives in the bridge/ layer because it spans the gap between the logic of the core and the wires of the hal.

package control

import (
	"context"
	"math"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bridge/hal"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/scheduler"
)

// Intent represents the high-level movement request from the AI.
// It is stored as a struct to allow for atomic swapping.
type Intent struct {
	Throttle float64 // 0.0 to 1.0
	Steer    float64 // -1.0 to 1.0
}

type ActuatorLoop struct {
	interlock    *hal.SafetyInterlock
	bus          *hal.BusProvider
	hz           int
	latestIntent atomic.Pointer[Intent] // Thread-safe storage for the latest AI command
}

// UpdateIntent is called by the Navigation/Kernel layer to feed new data to the bridge.
func (al *ActuatorLoop) UpdateIntent(newIntent *Intent) {
	al.latestIntent.Store(newIntent)
}

// fetchLatestIntent pulls the most recent command sent by the AI logic.
func (al *ActuatorLoop) fetchLatestIntent() Intent {
	ptr := al.latestIntent.Load()
	if ptr == nil {
		// Fallback to neutral if no intent has ever been sent
		return Intent{Throttle: 0.0, Steer: 0.0}
	}
	return *ptr
}

// translateIntentToHardware converts the abstract AI floats into a hardware-ready packet.
func (al *ActuatorLoop) translateIntentToHardware(intent Intent) hal.RawCommand {
	// 1. Clamp inputs to ensure they stay within safety bounds
	throttle := math.Max(0.0, math.Min(1.0, intent.Throttle))
	steer := math.Max(-1.0, math.Min(1.0, intent.Steer))

	// 2. Scale to hardware units (assuming 8-bit 0-255 range)
	// For steering, we center -1.0..1.0 around the neutral point 127
	return hal.RawCommand{
		ID: 0x01, // Main Actuator ID
		Data: []byte{
			byte(throttle * 255),          // 0 = Off, 255 = Full Power
			byte((steer + 1.0) * 127.5),   // 0 = Full Left, 127 = Center, 255 = Full Right
		},
	}
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
