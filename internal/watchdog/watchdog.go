//MIAUSEproject-founderKJ/multi-platform-AI/internal/watchdog/watchdog.go

package watchdog

import (
	"os"
	"sync"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type Config struct {
	TimeoutSeconds int
	OnFailure      string // "restart" | "degrade_to_safe_mode" | "halt"
}

type Watchdog struct {
	config   Config
	lastKick time.Time
	mu       sync.Mutex
	stopChan chan struct{}
}

func New(cfg Config) *Watchdog {
	return &Watchdog{
		config:   cfg,
		lastKick: time.Now(),
		stopChan: make(chan struct{}),
	}
}

// Start launches the monitoring loop
func (w *Watchdog) Start() {
	logging.Info("[WATCHDOG] Safety Interlock Active. Timeout: %ds", w.config.TimeoutSeconds)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				w.mu.Lock()
				elapsed := time.Since(w.lastKick)
				w.mu.Unlock()

				if elapsed.Seconds() > float64(w.config.TimeoutSeconds) {
					w.handleFailure()
				}
			case <-w.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Heartbeat (Kick) resets the timer. Call this in your main loops.
func (w *Watchdog) Heartbeat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastKick = time.Now()
}

func (w *Watchdog) handleFailure() {
	logging.Error("[WATCHDOG] CRITICAL FAILURE: Heartbeat Timeout. Execution Halted.")

	switch w.config.OnFailure {
	case "degrade_to_safe_mode":
		// Trigger hardware-level safe state (e.g., neutral gears, zero lasers)
		logging.Warn("[WATCHDOG] Transitioning to Safe Mode...")
		w.triggerSafeMode()
	case "restart":
		os.Exit(1) // Rely on OS/Systemd to restart
	default:
		os.Exit(1)
	}
}

func (w *Watchdog) triggerSafeMode() {
	// Implementation for physical safety:
	// 1. Cut power to non-essential buses
	// 2. Lock the vault
	// 3. Halt the node
	os.Exit(1)
}
