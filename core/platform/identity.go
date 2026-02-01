//MIAUSEproject-founderKJ/multi-platform-AI/core/platform/identity.go

package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"multi-platform-AI/configs/defaults"
	"multi-platform-AI/internal/logging"
)

// Finalize implements the scoring logic from the AIofSpeech reference
func (id *Identity) Finalize(env *defaults.EnvConfig) {
	logging.Info("[IDENTITY] Calculating Heuristic Platform Scores...")

	scores := make(map[PlatformClass]*PlatformScore)
	
	ensure := func(class PlatformClass, max float64) *PlatformScore {
		if _, ok := scores[class]; !ok {
			scores[class] = &PlatformScore{Class: class, MaxScore: max}
		}
		return scores[class]
	}

	// --- 1. Scoring Logic (Reference: InferPlatformClass) ---
	
	// VEHICLE Detection
	if env.Hardware.HasBus("can") {
		s := ensure(PlatformVehicle, 1.5)
		s.Score += 0.8 // Higher weight for physical CAN
		s.Signals = append(s.Signals, "Physical CAN-bus interface")
	}

	// ROBOTIC Detection
	if env.Hardware.HasBus("i2c") && env.Hardware.HasBus("spi") {
		s := ensure(PlatformRobot, 1.2)
		s.Score += 0.5
		s.Signals = append(s.Signals, "Micro-controller telemetry (I2C/SPI)")
	}

	// LAPTOP/MOBILE Detection (Sensors)
	if env.Hardware.HasSensor("battery") {
		s := ensure(PlatformLaptop, 1.0)
		s.Score += 0.6
		s.Signals = append(s.Signals, "Battery subsystem present")
	}

	// --- 2. Resolution Logic (Reference: ResolvePlatform) ---
	
	var bestClass PlatformClass = PlatformComputer
	var highConf float64 = 0.0
	var candidates []PlatformScore

	for _, s := range scores {
		s.Confidence = s.Score / s.MaxScore
		candidates = append(candidates, *s)
		
		if s.Confidence > highConf {
			highConf = s.Confidence
			bestClass = s.Class
		}
	}

	// Finalize Identity
	id.PlatformType = bestClass
	id.Source = "heuristic_scoring_v2"
	
	// Update the EnvConfig global state
	env.Platform.Final = bestClass
	env.Platform.Candidates = candidates
	env.Platform.ResolvedAt = time.Now()
	env.Platform.Locked = true

	logging.Info("[IDENTITY] Resolution: %s (Conf: %.2f)", bestClass, highConf)
	
	// --- 3. Attestation (Reference: FinalizeAttestation) ---
	generateHardwareHash(env)
}

func generateHardwareHash(env *defaults.EnvConfig) {
	// Unique string representing hardware reality
	rawState := fmt.Sprintf("%s-%s-%d",
		env.Identity.MachineID,
		env.Identity.OS,
		len(env.Hardware.Buses),
	)

	hash := sha256.Sum256([]byte(rawState))
	env.Attestation.EnvHash = hex.EncodeToString(hash[:])
	env.Attestation.Valid = true
	
	logging.Info("[SECURITY] Environment Hash Generated: %s...", env.Attestation.EnvHash[:12])
}

func (h *HardwareProfile) HasBus(target string) bool {
    for _, b := range h.Buses {
        if strings.EqualFold(b.Type, target) {
            return true
        }
    }
    return false
}