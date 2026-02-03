//MIAUSEproject-founderKJ/multi-platform-AI/internal/monitor/vitals.go
//To make the interface truly reflective, we need to monitor the "Biological" health of the hardware. The Vitals Monitor acts as the system's nervous system, gathering data on VRAM, CPU pressure, and thermal loads, then streaming them to your HMI (the "Face").

package monitor

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/configs/defaults"
)

type SystemVitals struct {
	CPULoad     float64 `json:"cpu_load"`
	MemoryUsage uint64  `json:"memory_usage"`
	VRAMUsage   uint64  `json:"vram_usage"`
	Temperature float64 `json:"temperature"`
	JitterMS    float64 `json:"jitter_ms"`
}

type VitalsMonitor struct {
	Env      *defaults.EnvConfig
	Ticker   *time.Ticker
	Stream   chan SystemVitals
}

func NewVitalsMonitor(env *defaults.EnvConfig) *VitalsMonitor {
	return &VitalsMonitor{
		Env:    env,
		Ticker: time.NewTicker(500 * time.Millisecond), // 2Hz Refresh
		Stream: make(chan SystemVitals, 10),
	}
}

func (m *VitalsMonitor) Start() {
	go func() {
		var lastTick time.Time
		for range m.Ticker.C {
			// Calculate Jitter (Deterministic Timing check)
			now := time.Now()
			jitter := 0.0
			if !lastTick.IsZero() {
				jitter = float64(now.Sub(lastTick).Milliseconds() - 500)
			}
			lastTick = now

			vitals := SystemVitals{
				CPULoad:     getCPULoad(),
				MemoryUsage: getRAMUsage(),
				VRAMUsage:   getVRAMUsage(m.Env.Hardware.Processors),
				Temperature: getPackageTemp(),
				JitterMS:    jitter,
			}
			m.Stream <- vitals
		}
	}()
}