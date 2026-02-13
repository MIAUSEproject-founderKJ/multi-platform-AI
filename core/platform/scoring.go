//MIAUSEproject-founderKJ/multi-platform-AI/core/platform/scoring.go
package platform

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
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
	var highestConf uint16

	logging.Info("[SCORING] Evaluating %d platform candidates...", len(env.Platform.Candidates))

	for i := range env.Platform.Candidates {
		c := &env.Platform.Candidates[i]

		// 1. Q16 Normalization Logic
		if c.MaxScore > 0 {
			// Avoid division by zero
			ratio := c.Score / c.MaxScore
			if ratio > 1.0 {
				ratio = 1.0 // Cap at 100%
			}
			c.Confidence = uint16(ratio * 65535)
		} else {
			c.Confidence = 0
		}

		// 2. Selection Logic (King of the Hill)
		if c.Confidence > highestConf {
			highestConf = c.Confidence
			bestCandidate = *c
		} else if c.Confidence == highestConf {
			// 3. Tie-Breaker: Safety Bias
			// If a Vehicle and a Workstation have the same score, assume Vehicle (Worst Case Safety)
			// This forces the system to load safety drivers, which is safer than omitting them.
			if isSafetyCritical(c.Type) && !isSafetyCritical(bestCandidate.Type) {
				bestCandidate = *c
				logging.Info("[SCORING] Tie-break won by safety-critical candidate: %s", c.Type)
			}
		}
	}

	// 4. Threshold Check
	// If the best we found is < 40% confidence, we don't trust it.
	const MinConfidenceThreshold = 26214 // ~40%
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