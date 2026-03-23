// boot/platform/identity_resolver.go
package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// RunPlatformInference scores and selects platform type based on EnvConfig and hardware fingerprint
func RunPlatformInference(env *schema.EnvConfig, fp probe.HardwareFingerprint) {

	logging.Info("[IDENTITY] Starting heuristic platform resolution...")

	scores := map[schema.PlatformClass]*schema.PlatformScore{}
	ensure := func(class schema.PlatformClass, max float64) *schema.PlatformScore {
		if scores[class] == nil {
			scores[class] = &schema.PlatformScore{Type: class, MaxScore: max}
		}
		return scores[class]
	}

	buses := env.Hardware.Buses
	osName := strings.ToLower(env.Identity.OS)

	// Vehicle / Robotic / Industrial scoring
	if hasBus(buses, "can") || osName == "qnx" || osName == "autosar" {
		s := ensure(schema.PlatformVehicle, 1.5)
		s.Score += 1.0
		s.Signals = append(s.Signals, "CAN bus / automotive RTOS detected")
	}
	if hasBus(buses, "i2c") && hasBus(buses, "spi") {
		s := ensure(schema.PlatformRobot, 1.2)
		s.Score += 0.4
		s.Signals = append(s.Signals, "I2C+SPI sensors detected")
	}

	// Desktop/Laptop scoring
	collectDesktopSignals(env, fp, scores)

	// Compute final confidence
	var best schema.PlatformClass = schema.PlatformComputer
	var highConf mathutil.Q16
	for _, s := range scores {
		s.Confidence = mathutil.Q16(mathutil.FromFloat64(s.Score / s.MaxScore))
		if s.Confidence > highConf {
			highConf = s.Confidence
			best = s.Type
		}
	}

	env.Platform.Final = best
	env.Platform.Locked = true
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[DEBUG] Identity Resolver RunPlatformInference: best=%s, highConf=%d%%", best, highConf.Percentage())
	logging.Info("[IDENTITY] Identity Resolver Resolution: %s (Conf: %d%%)", best, highConf.Percentage())
}

// ------------------- helpers -------------------

func hasBus(buses []schema.BusCapability, busType string) bool {
	for _, b := range buses {
		if strings.EqualFold(b.Type, busType) {
			return true
		}
	}
	return false
}

// Desktop/Laptop scoring using fingerprint
func collectDesktopSignals(env *schema.EnvConfig, fp probe.HardwareFingerprint, scores map[schema.PlatformClass]*schema.PlatformScore) {
	s := scores[schema.PlatformComputer]
	if s == nil {
		s = &schema.PlatformScore{Type: schema.PlatformComputer, MaxScore: 1.5}
		scores[schema.PlatformComputer] = s
	}

	cpu := runtime.NumCPU()
	scoreFromCPU := 0.1 * float64(cpu)
	scoreFromPCI := 0.1 * float64(len(fp.PCI))
	scoreFromMAC := 0.1 * float64(len(fp.MAC))
	scoreFromBattery := 0.0
	if env.Hardware.HasBattery {
		scoreFromBattery = 0.3
	}

	s.Score += scoreFromCPU + scoreFromPCI + scoreFromMAC + scoreFromBattery
	s.Signals = append(s.Signals,
		fmt.Sprintf("CPU cores: %d", cpu),
		fmt.Sprintf("Battery: %v", env.Hardware.HasBattery),
		fmt.Sprintf("PCI devices: %d", len(fp.PCI)),
		fmt.Sprintf("MAC addresses: %d", len(fp.MAC)),
	)

	// Update Candidates slice for logging
	env.Platform.Candidates = append(env.Platform.Candidates, *s)
}

// ComputeHardwareFingerprint remains unchanged
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
