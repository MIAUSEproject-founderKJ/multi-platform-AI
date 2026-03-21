// boot/platform/identity_resolver.go
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

// -----------------------------
// Platform Inference & Scoring
// -----------------------------

func RunPlatformInference(env *schema.EnvConfig) {
	logging.Info("[IDENTITY] Starting heuristic platform resolution...")

	osName := strings.ToLower(env.Identity.OS)
	buses := env.Hardware.Buses
	scores := map[schema.PlatformClass]*schema.PlatformScore{}

	ensure := func(class schema.PlatformClass, max float64) *schema.PlatformScore {
		if scores[class] == nil {
			scores[class] = &schema.PlatformScore{Type: class, MaxScore: max}
		}
		return scores[class]
	}

	// Vehicle detection
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

	// Robotics / Industrial
	if hasBus(buses, "i2c") && hasBus(buses, "spi") {
		s := ensure(schema.PlatformRobot, 1.2)
		s.Score += 0.4
		s.Signals = append(s.Signals, "I2C+SPI sensors")
	}

	// Laptop / Mobile
	if env.Hardware.HasBattery {
		s := ensure(schema.PlatformLaptop, 1.0)
		s.Score += 0.6
		s.Signals = append(s.Signals, "Battery present")
	} else {
		s := ensure(schema.PlatformComputer, 1.0)
		s.Score += 0.4
	}

	// Convert scores to Q16 confidence and select the best platform
	var bestClass schema.PlatformClass = schema.PlatformComputer
	var highConf mathutil.Q16

	var candidates []schema.PlatformScore

	for _, s := range scores {
		s.Confidence = mathutil.Q16(mathutil.FromFloat64(s.Score / s.MaxScore))
		candidates = append(candidates, *s)
		if s.Confidence > highConf {
			highConf = s.Confidence
			bestClass = s.Type
		}
	}

	env.Platform.Candidates = candidates
	env.Platform.Final = bestClass
	env.Platform.Locked = true
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[IDENTITY] Resolution: %s (Conf: %d%%)", bestClass, highConf.Percentage())

	ComputeHardwareFingerprint(env)
}

// -----------------------------
// Helpers
// -----------------------------

func hasBus(buses []schema.BusCapability, busType string) bool {
	for _, b := range buses {
		if strings.EqualFold(b.Type, busType) {
			return true
		}
	}
	return false
}

// -----------------------------
// Attestation / Fingerprint
// -----------------------------

func ComputeHardwareFingerprint(env *schema.EnvConfig) []byte {
	payload := fmt.Sprintf("%s|%s|%s|%d|%d",
		env.Identity.MachineID,
		env.Identity.OS,
		env.Identity.Arch,
		len(env.Hardware.Processors),
		len(env.Hardware.Buses),
	)

	hash := sha256.Sum256([]byte(payload))
	env.Attestation.EnvHash = hex.EncodeToString(hash[:])
	env.Attestation.Valid = true

	logging.Info("[SECURITY] Environment Hash Generated: %s...", env.Attestation.EnvHash[:12])
	return hash[:]
}

// -----------------------------
// Security Profile Classification
// -----------------------------

type SecurityProfile struct {
	PlatformClass        string
	RequiresBiometry     bool
	RequiresKeyTelemetry bool
	AuthTimeoutMinutes   int
}
