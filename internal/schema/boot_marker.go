//internal/schema/boot_marker.go

package schema

import "time"

type FirstBootMarker struct {
	MachineName   string     `json:"machine_name"`
	SchemaVersion int        `json:"schema_version"`
	GoldenHash    []byte     `json:"golden_hash"`
	Initialized   bool       `json:"initialized"`
	CreatedAt     time.Time  `json:"created_at"`
	TrustLevel    TrustLevel `json:"trust_level"`
}
