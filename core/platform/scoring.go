//MIAUSEproject-founderKJ/multi-platform-AI/core/platform/scoring.go

package platform

import "time"

type PlatformScore struct {
	Class      PlatformClass
	Score      float64
	MaxScore   float64
	Confidence float64
	Signals    []string
}

type PlatformDescriptor struct {
	Final      PlatformClass
	Candidates []PlatformScore
	Source     string
	ResolvedAt time.Time
	Locked     bool
}