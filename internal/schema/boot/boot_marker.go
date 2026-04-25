//internal\schema\boot\boot_marker.go

package schema_boot

import (
	"time"

	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

type FirstBootMarker struct {
	MachineID     string                  `json:"machine_id"`
	SchemaVersion int                     `json:"schema_version"`
	GoldenHash    []byte                  `json:"golden_hash"`
	Initialized   bool                    `json:"initialized"`
	CreatedAt     time.Time               `json:"created_at"`
	BootTrust     schema_system.BootTrust `json:"boot_trust"`
}
