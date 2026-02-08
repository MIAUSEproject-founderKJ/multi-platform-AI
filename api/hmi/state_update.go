//api/hmi/state_update.go

package hmi

import (
	"time"
)

// StateUpdate is the high-frequency telemetry packet sent to the HUD.
type StateUpdate struct {
	// Timing
	Timestamp time.Time `json:"timestamp"`

	// Bayesian Identity
	TrustScore  float64 `json:"trust_score"`   // 0.0 to 1.0 (for UI Gauges)
	RawScoreQ16 uint16  `json:"raw_score_q16"` // 0 to 65535 (for System Logs)

	// Operational State
	Mode        string `json:"mode"`         // AUTONOMOUS, ASSISTED, MANUAL_ONLY
	StatusLabel string `json:"status_label"` // NOMINAL, DEGRADED, CRITICAL
	
	// Platform Data
	PlatformClass string `json:"platform_class"`
	ActiveBuses   int    `json:"active_buses"`
	
	// Message/Alert
	LastMessage string `json:"last_message,omitempty"`
}

// ProgressUpdate is used specifically during the Boot Sequence (Stage 1).
type ProgressUpdate struct {
	Percentage float32 `json:"percentage"`
	Message    string  `json:"message"`
	Stage      string  `json:"stage"` // e.g., "BOOT_PROBE", "ATTESTATION"
}