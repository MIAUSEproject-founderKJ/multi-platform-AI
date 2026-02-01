//MIAUSEproject-founderKJ/multi-platform-AI/configs/platforms/resolution.go

package platform

import (
    "multi-platform-AI/configs/platforms"
    "time"
)

type PlatformScore struct {
    Class      platforms.PlatformClass `json:"class"`
    Score      float64                `json:"score"`
    Confidence uint16                 `json:"confidence"`
    Signals    []string               `json:"signals"`
}

type PlatformResolution struct {
    Candidates []PlatformScore         `json:"candidates"`
    Final      platforms.PlatformClass `json:"final"`
    Locked     bool                    `json:"locked"`
    Source     string                  `json:"source"`
    ResolvedAt time.Time               `json:"resolved_at"`
}