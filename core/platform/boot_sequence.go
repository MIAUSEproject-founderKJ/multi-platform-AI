//MIAUSEproject-founderKJ/multi-platform-AI/core/platform/boot_sequence.go

package platform

import (
	"fmt"

	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// Summary returns a human-readable log line for the console.
func (bs *BootSequence) Summary() string {
	return fmt.Sprintf("[%s] Booted as %s | Trust: %.0f%% | Mode: %s",
		bs.Timestamp.Format("15:04:05"),
		bs.PlatformID,
		bs.TrustScore*100,
		bs.Mode,
	)
}

// CanOperate returns true if the system allows any form of actuation.
func (bs *BootSequence) CanOperate() bool {
	return bs.Mode != "MANUAL_ONLY"
}

// IsAutonomous returns true only if the trust is high enough for self-governance.
func (bs *BootSequence) IsAutonomous() bool {
	return bs.Mode == "AUTONOMOUS"
}
