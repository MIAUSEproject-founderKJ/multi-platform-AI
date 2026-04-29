// bootstrap/platform/resolve.go
package platform

import (
	"time"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

// ResolvePlatform selects the final operational class based on scores and attestation locks.
func ResolvePlatform(env *internal_environment.EnvConfig) internal_environment.PlatformClass {
	// 1. Check for Attestation Lock (The "Immutable" path)
	if env.Attestation.Locked {
		env.Platform.Final = env.Attestation.PlatformClass
		env.Platform.Locked = true
		env.Platform.Source = "attestation_lock"
		return env.Platform.Final
	}

	// 2. Fallback to Score-based Resolution
	var bestType internal_environment.PlatformClass
	var highestScore float64 = -1.0

	for _, candidate := range env.Platform.Candidates {
		if candidate.Score > highestScore {
			highestScore = candidate.Score
			bestType = candidate.Type
		}
	}

	env.Platform.Final = bestType
	env.Platform.ResolvedAt = time.Now()
	env.Platform.Source = "probabilistic_match"

	return bestType
}
