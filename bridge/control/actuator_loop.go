//MIAUSEproject-founderKJ/multi-platform-AI/bridge/control/actuator_loop.go

//This file handles the high-frequency translation of AI "Intents" into hardware "Commands." It lives in the bridge/ layer because it spans the gap between the logic of the core and the wires of the hal.

package control

import (
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bridge/hal"
)

type Intent struct {
	Throttle float64
	Steer    float64
}

type ActuatorLoop struct {
	interlock    *hal.SafetyInterlock
	bus          hal.BusProvider // Interface is better than pointer to struct here
	hz           int
	latestIntent atomic.Pointer[Intent]
}

// ... fetchLatestIntent and translateIntentToHardware remain largely the same ...

func (al *ActuatorLoop) tick() {
	if al.bus == nil || al.interlock == nil {
		return
	}

	// 1. GATING: Use atomic to check safety state
	// We cast the result to InterlockState for readability
	currentState := hal.InterlockState(atomic.LoadInt32(&al.interlock.currentState))

	if currentState != hal.StateClear {
		// Hardware safety trigger (e.g., physical E-Stop)
		al.bus.WriteSafeState()
		return
	}

	// 2. FETCH & TRANSLATE
	intent := al.fetchLatestIntent()
	rawCommand := al.translateIntentToHardware(intent)

	// 3. EXECUTE: The scaling formula used here is:
	// $$byte = (steer + 1.0) \times 127.5$$
	al.bus.Write(rawCommand)
}
