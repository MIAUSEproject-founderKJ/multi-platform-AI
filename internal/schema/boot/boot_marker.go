//internal/schema/boot/boot_marker.go

package internal_boot

import (
	"time"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

type FirstBootMarker struct {
	MachineID     string                         `json:"machine_id"`
	SchemaVersion int                            `json:"internal_version"`
	GoldenHash    []byte                         `json:"golden_hash"`
	Initialized   bool                           `json:"initialized"`
	CreatedAt     time.Time                      `json:"created_at"`
	BootTrust     internal_environment.BootTrust `json:"boot_trust"`
}
