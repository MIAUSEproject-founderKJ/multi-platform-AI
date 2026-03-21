// boot/platform/scoring.go
package platform

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// Define constants to prevent magic strings
const (
	ClassUnknown     schema.PlatformClass = "Unknown"
	ClassVehicle     schema.PlatformClass = "Automotive"
	ClassIndustrial  schema.PlatformClass = "Industrial"
	ClassWorkstation schema.PlatformClass = "Workstation"
)

// RunResolution determines the Final PlatformClass based on confidence scores.
func RunResolution(env *schema.EnvConfig) {
	// Security: If the platform is locked (e.g., by a hardware dongle), do not override.
	if env.Platform.Locked {
		return
	}

	var bestCandidate schema.PlatformScore
	var highestConf mathutil.Q16

	for i := range env.Platform.Candidates {
		c := &env.Platform.Candidates[i]

		// Normalize score → Q16
		if c.MaxScore > 0 {
			ratio := c.Score / c.MaxScore
			if ratio > 1.0 {
				ratio = 1.0
			}
			c.Confidence = mathutil.Q16(mathutil.FromFloat64(ratio))
		} else {
			c.Confidence = 0
		}

		// Selection
		if c.Confidence > highestConf {
			highestConf = c.Confidence
			bestCandidate = *c
		} else if c.Confidence == highestConf {
			if isSafetyCritical(c.Type) && !isSafetyCritical(bestCandidate.Type) {
				bestCandidate = *c
			}
		}
	}

	// 4. Threshold Check
	// If the best we found is < 40% confidence, we don't trust it.
	const MinConfidenceThreshold = mathutil.Q16(26214) // 40% of 65535
	if highestConf < MinConfidenceThreshold {
		logging.Warn("[SCORING] Ambiguous hardware identity (Conf: %d). Defaulting to SAFE_MODE.", highestConf)
		env.Platform.Final = "Generic_SafeMode"
	} else {
		env.Platform.Final = bestCandidate.Type
	}

	env.Platform.Source = "heuristic_v1"
	env.Platform.ResolvedAt = time.Now()
	env.Platform.Locked = true

	logging.Info("[SCORING] Resolution Complete. Identity: %s (Confidence: %d/65535)",
		env.Platform.Final, highestConf)
}

// Helper to prioritize safety-critical platforms during ties
func isSafetyCritical(pType schema.PlatformClass) bool {
	return pType == ClassVehicle || pType == ClassIndustrial
}
