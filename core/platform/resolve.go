// MIAUSEproject-founderKJ/multi-platform-AI/core/platform/resolve.go
package platform

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ResolvePlatform selects the final operational class based on scores and attestation locks.
func ResolvePlatform(env *schema.EnvConfig) schema.PlatformClass {
	// 1. Check for Attestation Lock (The "Immutable" path)
	if env.Attestation.schema.Locked {
		env.Platform.Final = env.Attestation.Platform
		env.Platform.Locked = true
		env.Platform.Source = "attestation_lock"
		return env.Platform.Final
	}

	// 2. Fallback to Score-based Resolution
	var bestClass schema.PlatformClass
	var highestScore float64 = -1.0

	for _, candidate := range env.Platform.Candidates {
		if candidate.Score > highestScore {
			highestScore = candidate.Score
			bestClass = candidate.Class
		}
	}

	env.Platform.Final = bestClass
	env.Platform.ResolvedAt = time.Now()
	env.Platform.Source = "probabilistic_match"

	return bestClass
}
