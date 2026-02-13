//MIAUSEproject-founderKJ/multi-platform-AI/core/platform/identity_resolver.go

package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ResolveIdentity performs the full pipeline: Scoring -> Resolution -> Attestation
func ResolveIdentity(env *schema.EnvConfig) {
	logging.Info("[IDENTITY] Starting heuristic platform resolution...")

	// 1. Run Inference (Reference: InferPlatformClass)
	scores := runPlatformInference(env)

	// 2. Resolve Final Class (Reference: ResolvePlatform)
	finalizePlatform(env, scores)

	// 3. Generate Hardware Fingerprint (Reference: FinalizeAttestation)
	performAttestation(env)
}

// runPlatformInference calculates the probability of each platform class
func runPlatformInference(env *schema.EnvConfig) []schema.PlatformScore {
	osName := strings.ToLower(env.Identity.OS)
	buses := env.Hardware.Buses
	scores := map[schema.PlatformClass]*schema.PlatformScore{}

	ensure := func(class schema.PlatformClass, max float64) *schema.PlatformScore {
		if scores[class] == nil {
			scores[class] = &schema.PlatformScore{Class: class, MaxScore: max}
		}
		return scores[class]
	}

	// VEHICLE Logic
	if hasBus(buses, "can") {
		s := ensure(schema.PlatformVehicle, 1.5)
		s.Score += 0.5
		s.Signals = append(s.Signals, "CAN bus detected")
	}
	if osName == "qnx" || osName == "autosar" {
		s := ensure(schema.PlatformVehicle, 1.5)
		s.Score += 1.0
		s.Signals = append(s.Signals, "Automotive RTOS")
	}

	// ROBOTIC Logic
	if hasBus(buses, "i2c") && hasBus(buses, "spi") {
		s := ensure(schema.PlatformRobot, 1.2)
		s.Score += 0.4
		s.Signals = append(s.Signals, "I2C+SPI sensors")
	}

	// Normalize into Q16 Confidence (matching your env_config.go types)
	var out []schema.PlatformScore
	for _, s := range scores {
		if s.MaxScore > 0 {
			s.Confidence = mathutil.FromFloat64(s.Score / s.MaxScore)
		}
		out = append(out, *s)
	}
	return out
}

// finalizePlatform selects the "Best Fit" and locks the state
func finalizePlatform(env *schema.EnvConfig, scores []schema.PlatformScore) {
	if len(scores) == 0 {
		env.Platform.Final = schema.PlatformComputer
		env.Platform.Source = "fallback_default"
	} else {
		var best schema.PlatformScore
		for _, s := range scores {
			if s.Confidence.ToFloat64() > best.Confidence.ToFloat64() {
				best = s
			}
		}
		env.Platform.Candidates = scores
		env.Platform.Final = best.Class
		env.Platform.Source = "heuristic_scoring_v1"
	}
	env.Platform.ResolvedAt = time.Now()
	env.Platform.Locked = true
}

// performAttestation creates the crypto-link between code and hardware
func performAttestation(env *schema.EnvConfig) {
	rawState := fmt.Sprintf("%s-%s-%d",
		env.Identity.MachineName,
		env.Identity.OS,
		len(env.Hardware.Buses),
	)
	hash := sha256.Sum256([]byte(rawState))
	env.Attestation.EnvHash = hex.EncodeToString(hash[:])

	// Logic: If we are on a platform with a hardware TPM, it's Strong.
	// If we are just hashing a string in software, it's actually Weak.
	if hasSecureEnclave() {
		env.Attestation.Level = schema.AttestationStrong // Score will be 0.99
	} else {
		env.Attestation.Level = schema.AttestationWeak // Score will be 0.75
	}
	logging.Info("[SECURITY] Attestation Hash: %s", env.Attestation.EnvHash[:12])
	env.Attestation.Valid = true
}

func hasBus(buses []schema.BusCapability, busType string) bool {
	for _, b := range buses {
		if strings.EqualFold(b.Type, busType) {
			return true
		}
	}
	return false
}
