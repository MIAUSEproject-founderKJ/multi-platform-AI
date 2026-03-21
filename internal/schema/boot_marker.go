//internal/schema/boot_marker.go

package schema

import "time"

type FirstBootMarker struct {
	MachineID     string    `json:"machine_id"`
	SchemaVersion int       `json:"schema_version"`
	GoldenHash    []byte    `json:"golden_hash"`
	Initialized   bool      `json:"initialized"`
	CreatedAt     time.Time `json:"created_at"`
	BootTrust     BootTrust `json:"boot_trust"`
}
