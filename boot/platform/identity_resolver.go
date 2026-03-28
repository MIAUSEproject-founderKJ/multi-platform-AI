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

	osName := strings.ToLower(env.Identity.OS)

	// Vehicle / Robotic / Industrial scoring
	if hasBus(fp, "can") || osName == "qnx" || osName == "autosar" {
		s := ensure(schema.PlatformVehicle, 1.5)
		s.Score += 1.0
		s.Signals = append(s.Signals, schema.Signal{
			Name:       "can_bus_or_rtos",
			Value:      1.0,
			Weight:     1.0,
			Confidence: 0.9,
			Source:     "heuristic",
		})
	}
	if hasBus(fp, "i2c") && hasBus(fp, "spi") {
		s := ensure(schema.PlatformRobot, 1.2)
		s.Score += 0.4

		s.Signals = append(s.Signals, schema.Signal{
			Name:       "I2C+SPI sensors",
			Value:      1.0,
			Weight:     1.0,
			Confidence: 0.9,
			Source:     "heuristic",
		})
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
func hasBus(fp probe.HardwareFingerprint, bus string) bool {
	return fp.Buses[bus]
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
		schema.Signal{
			Name:       "cpu_cores",
			Value:      minFloat(float64(cpu)/16.0, 1.0),
			Weight:     0.3,
			Confidence: 0.9,
			Source:     "runtime",
		},
		schema.Signal{
			Name:       "battery_present",
			Value:      boolToFloat(env.Hardware.HasBattery),
			Weight:     0.3,
			Confidence: 0.95,
			Source:     "power",
		},
		schema.Signal{
			Name:       "pci_devices",
			Value:      minFloat(float64(len(fp.PCI))/10.0, 1.0),
			Weight:     0.3,
			Confidence: 0.9,
			Source:     "runtime",
		},
		schema.Signal{
			Name:       "mac_addresses",
			Value:      minFloat(float64(len(fp.MAC))/5.0, 1.0),
			Weight:     0.2,
			Confidence: 0.9,
			Source:     "runtime",
		},
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

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
