//MIAUSEproject-founderKJ/multi-platform-AI/core/platform/identity_resolver.go

package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"multi-platform-AI/configs/configStruct"
	"multi-platform-AI/internal/logging"
)

// ResolveIdentity performs the full pipeline: Scoring -> Resolution -> Attestation
func ResolveIdentity(env *configStruct.EnvConfig) {
	logging.Info("[IDENTITY] Starting heuristic platform resolution...")

	// 1. Run Inference (Reference: InferPlatformClass)
	scores := runPlatformInference(env)

	// 2. Resolve Final Class (Reference: ResolvePlatform)
	finalizePlatform(env, scores)

	// 3. Generate Hardware Fingerprint (Reference: FinalizeAttestation)
	performAttestation(env)
}

// runPlatformInference calculates the probability of each platform class
func runPlatformInference(env *configStruct.EnvConfig) []configStruct.PlatformScore {
	osName := strings.ToLower(env.Identity.OS)
	buses := env.Hardware.Buses
	scores := map[configStruct.PlatformClass]*configStruct.PlatformScore{}

	ensure := func(class configStruct.PlatformClass, max float64) *configStruct.PlatformScore {
		if scores[class] == nil {
			scores[class] = &configStruct.PlatformScore{Class: class, MaxScore: max}
		}
		return scores[class]
	}

	// VEHICLE Logic
	if hasBus(buses, "can") {
		s := ensure(configStruct.PlatformVehicle, 1.5)
		s.Score += 0.5
		s.Signals = append(s.Signals, "CAN bus detected")
	}
	if osName == "qnx" || osName == "autosar" {
		s := ensure(configStruct.PlatformVehicle, 1.5)
		s.Score += 1.0
		s.Signals = append(s.Signals, "Automotive RTOS")
	}

	// ROBOTIC Logic
	if hasBus(buses, "i2c") && hasBus(buses, "spi") {
		s := ensure(configStruct.PlatformRobot, 1.2)
		s.Score += 0.4
		s.Signals = append(s.Signals, "I2C+SPI sensors")
	}

	// Normalize into Q16 Confidence (matching your env_config.go types)
	var out []configStruct.PlatformScore
	for _, s := range scores {
		if s.MaxScore > 0 {
			s.Confidence = configStruct.Q16FromFloat(s.Score / s.MaxScore)
		}
		out = append(out, *s)
	}
	return out
}

// finalizePlatform selects the "Best Fit" and locks the state
func finalizePlatform(env *configStruct.EnvConfig, scores []configStruct.PlatformScore) {
	if len(scores) == 0 {
		env.Platform.Final = configStruct.PlatformComputer
		env.Platform.Source = "fallback_default"
	} else {
		var best configStruct.PlatformScore
		for _, s := range scores {
			if s.Confidence.Float64() > best.Confidence.Float64() {
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
func performAttestation(env *configStruct.EnvConfig) {
	// Refined rawState to use the unified MachineID
	rawState := fmt.Sprintf("%s-%s-%d",
		env.Identity.MachineName, // Mapping from MachineIdentity
		env.Identity.OS,
		len(env.Hardware.Buses),
	)

	hash := sha256.Sum256([]byte(rawState))
	env.Attestation.EnvHash = hex.EncodeToString(hash[:])
	env.Attestation.Valid = true
	env.Attestation.Level = configStruct.AttestationStrong
	
	logging.Info("[SECURITY] Attestation Hash: %s", env.Attestation.EnvHash[:12])
}

func hasBus(buses []configStruct.BusCapability, busType string) bool {
	for _, b := range buses {
		if strings.EqualFold(b.Type, busType) {
			return true
		}
	}
	return false
}