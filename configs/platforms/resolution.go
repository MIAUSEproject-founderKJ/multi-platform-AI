//MIAUSEproject-founderKJ/multi-platform-AI/configs/platforms/resolution.go

package platforms

import (
	"time"
)

// PlatformScore tracks the heuristic weight for a specific platform type.
type PlatformScore struct {
	Class      PlatformClass `json:"class"`
	Score      float64       `json:"score"`      // Raw cumulative score
	MaxScore   float64       `json:"max_score"`  // Potential maximum for normalization
	Confidence uint16        `json:"confidence"` // Normalized Q16 (0-65535)
	Signals    []string      `json:"signals"`    // Evidence found (e.g., "CAN_BUS_PRESENT")
}

// PlatformResolution is the finalized identity of the environment.
type PlatformResolution struct {
	Candidates []PlatformScore `json:"candidates"`
	Final      PlatformClass   `json:"final"`
	Locked     bool            `json:"locked"`
	Source     string          `json:"source"` // e.g., "heuristic_v1" or "manual_override"
	ResolvedAt time.Time       `json:"resolved_at"`
}