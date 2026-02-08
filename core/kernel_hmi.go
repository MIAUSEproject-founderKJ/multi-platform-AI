//core/kernel_hmi.go

package core

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/hmi"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

// RunHMILoop manages the lifecycle of the Perception Overlay and Telemetry HUD.
func (k *Kernel) RunHMILoop() {
	logging.Info("[HMI] Perception Loop Started (30Hz Telemetry)")

	// Create a ticker to maintain a steady UI framerate (approx 33ms = 30fps)
	ticker := time.NewTicker(33 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-k.ctx.Done():
			logging.Info("[HMI] Terminating HMI Loop (Context Closed)")
			return
		case <-ticker.C:
			// 1. Snapshot the current Bayesian Trust State
			trust := k.Evaluator.Evaluate(k.EnvConfig)

			// 2. Construct the Telemetry Packet
			update := hmi.StateUpdate{
				Timestamp:     time.Now(),
				TrustScore:    trust.CurrentScore,
				RawScoreQ16:   trust.RawScoreQ16,
				Mode:          trust.OperationMode,
				StatusLabel:   trust.Label,
				PlatformClass: string(k.EnvConfig.Platform.Final),
			}

			// 3. Push to the HMI Pipe
			// This feeds the VisionHUD we linked in the BootManager
			select {
			case k.HMIPipe <- update:
				// Successfully sent to HUD
			default:
				// Pipe is full; skip this frame to prevent blocking the Kernel
			}
		}
	}
}
