//MIAUSEproject-founderKJ/multi-platform-AI/core/platform/identity.go

package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type Identity struct {
	PlatformType schema.PlatformClass
	Source       string
}

// Finalize implements the scoring logic using the unified schema.
func (id *Identity) Finalize(env *schema.EnvConfig) {
	logging.Info("[IDENTITY] Calculating Heuristic Platform Scores...")

	// Use the types directly from schema
	scores := make(map[schema.PlatformClass]*schema.PlatformScore)

	ensure := func(class schema.PlatformClass, max float64) *schema.PlatformScore {
		if _, ok := scores[class]; !ok {
			scores[class] = &schema.PlatformScore{Class: class, MaxScore: max}
		}
		return scores[class]
	}

	// --- 1. Scoring Logic ---

	// Vehicle Detection
if hasBus(env.Hardware, "can") {
    s := ensure(schema.PlatformVehicle, 1.5)
    s.Score += 0.8
    s.Signals = append(s.Signals, "CAN bus detected")
}

// Robot/Industrial Detection
if hasBus(env.Hardware, "i2c") && hasBus(env.Hardware, "spi") {
    s := ensure(schema.PlatformRobot, 1.2)
    s.Score += 0.5
    s.Signals = append(s.Signals, "I2C/SPI detected")
}

	// ROBOTIC Detection
	if hasBus(env.Hardware, "i2c") && hasBus(env.Hardware, "spi") {
		s := ensure(schema.PlatformRobot, 1.2)
		s.Score += 0.5
		s.Signals = append(s.Signals, "Micro-controller telemetry (I2C/SPI)")
	}

	// LAPTOP/MOBILE Detection
	// Laptop/Desktop Detection
if env.Hardware.HasBattery {
    s := ensure(schema.PlatformLaptop, 1.0)
    s.Score += 0.6
    s.Signals = append(s.Signals, "Battery present")
} else {
    s := ensure(schema.PlatformComputer, 1.0)
    s.Score += 0.4
}

	// --- 2. Resolution Logic ---

	var bestClass schema.PlatformClass = schema.PlatformComputer
	var highConf float64 = 0.0
	var candidates []schema.PlatformScore

	for _, s := range scores {
		// Inside identity.go
		s.Confidence = uint16((s.Score / s.MaxScore) * 65535)
		candidates = append(candidates, *s)

		if s.Confidence > highConf {
			highConf = float64(s.Confidence)
			bestClass = s.Class
		}
	}

	id.PlatformType = bestClass
	id.Source = "heuristic_scoring_v2"

	// Update the EnvConfig global state
	env.Platform.Final = bestClass
	env.Platform.Candidates = candidates
	env.Platform.ResolvedAt = time.Now()
	env.Platform.Locked = true

	logging.Info("[IDENTITY] Resolution: %s (Conf: %.2f)", bestClass, highConf)

	generateHardwareHash(env)
}

// hasBus is now a helper function because we can't add methods to schema.HardwareProfile
func hasBus(h schema.HardwareProfile, target string) bool {
	for _, b := range h.Buses {
		if strings.EqualFold(b.Type, target) {
			return true
		}
	}
	return false
}

func generateHardwareHash(env *schema.EnvConfig) {
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

// DetectPlatformClass performs the initial "Discovery" phase
func DetectPlatformClass(hw *schema.HardwareProfile) schema.PlatformClass {
	// 1. Check for Vehicle Indicators (CAN-bus)
	if _, err := os.Stat("/sys/class/net/can0"); err == nil {
		hw.Buses = append(hw.Buses, schema.BusCapability{
			ID:         "can0",
			Type:       "can",
			Confidence: 65535,
			Source:     "probed",
		})
		return schema.PlatformVehicle
	}

	// 2. Check for Industrial Indicators
	if os.Getenv("INDUSTRIAL_NODE_ID") != "" {
		return schema.PlatformIndustrial
	}

	// 3. Check for Robotics/Embedded Indicators
	if _, err := os.Stat("/dev/i2c-1"); err == nil {
		hw.Buses = append(hw.Buses, schema.BusCapability{
			ID:         "i2c_bus_1",
			Type:       "i2c",
			Confidence: 65535,
			Source:     "probed",
		})
		return schema.PlatformRobot
	}

	return schema.PlatformComputer
}

func InitializeIdentity() schema.MachineIdentity {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "unknown-node"
	}

	return schema.MachineIdentity{
		MachineName: hostname,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
	}
}
