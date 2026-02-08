//MIAUSEproject-founderKJ/multi-platform-AI/internal/monitor/vitals.go
//To make the interface truly reflective, we need to monitor the "Biological" health of the hardware. The Vitals Monitor acts as the system's nervous system, gathering data on VRAM, CPU pressure, and thermal loads, then streaming them to your HMI (the "Face").

package monitor

import (
	"runtime"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema" // Fixed: Use schema to break import cycle
)

type SystemVitals struct {
	CPULoad     float64 `json:"cpu_load"`
	MemoryUsage uint64  `json:"memory_usage"`
	VRAMUsage   uint64  `json:"vram_usage"`
	Temperature float64 `json:"temperature"`
	JitterMS    float64 `json:"jitter_ms"`
}

type VitalsMonitor struct {
	Env    *schema.EnvConfig
	Ticker *time.Ticker
	Stream chan SystemVitals
}

func NewVitalsMonitor(env *schema.EnvConfig) *VitalsMonitor {
	return &VitalsMonitor{
		Env:    env,
		Ticker: time.NewTicker(500 * time.Millisecond),
		Stream: make(chan SystemVitals, 10),
	}
}

func (m *VitalsMonitor) Start() {
	go func() {
		var lastTick time.Time
		for range m.Ticker.C {
			now := time.Now()
			jitter := 0.0
			if !lastTick.IsZero() {
				// Jitter = Actual time - Expected interval (500ms)
				jitter = float64(now.Sub(lastTick).Milliseconds() - 500)
			}
			lastTick = now

			m.Stream <- SystemVitals{
				CPULoad:     getCPULoad(),
				MemoryUsage: getRAMUsage(),
				VRAMUsage:   getVRAMUsage(),
				Temperature: getPackageTemp(),
				JitterMS:    jitter,
			}
		}
	}()
}

// --- Hardware Probing Helpers ---

func getCPULoad() float64 {
	// In a real implementation, use gopsutil or read /proc/stat
	return 0.0
}

func getRAMUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func getVRAMUsage() uint64 {
	// Logic to probe NVIDIA/Intel/Apple Silicon VRAM
	return 0
}

func getPackageTemp() float64 {
	// Default to a safe 'simulated' value
	return 42.0
}
