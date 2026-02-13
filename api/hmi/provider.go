//api/hmi/provider.go

package hmi

import "time"

// SystemPulse represents the health telemetry of the node.
type SystemPulse struct {
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"mem_usage"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

// MonitorProgress for long-running AI tasks (like distillation).
type MonitorProgress struct {
	Task     string  `json:"task"`
	Progress float64 `json:"progress"` // 0.0 to 1.0
}

// HMIPipe defines how the core sends data to the User Interface.
type HMIPipe interface {
	SendTelemetry(pulse SystemPulse)
	SendSpatial(frame interface{}) // For Voxel/Lidar data
	SendProgress(update MonitorProgress)
}
