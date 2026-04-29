// bootstrap/platform/identity_resolver.go
package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/math_convert"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

// Desktop/Laptop scoring using fingerprint
func collectDesktopSignals(env *internal_environment.EnvConfig, fp probe.HardwareFingerprint, scores map[internal_environment.PlatformClass]*internal_environment.PlatformScore) {
	s := scores[internal_environment.PlatformComputer]
	if s == nil {
		s = &internal_environment.PlatformScore{Type: internal_environment.PlatformComputer, MaxScore: 1.5}
		scores[internal_environment.PlatformComputer] = s
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
		internal_environment.Signal{
			Name:       "cpu_cores",
			Value:      minFloat(float64(cpu)/16.0, 1.0),
			Weight:     0.3,
			Confidence: math_convert.FromFloat64(0.9),
			Source:     "runtime",
		},
		internal_environment.Signal{
			Name:       "battery_present",
			Value:      BoolToFloat(env.Hardware.HasBattery),
			Weight:     0.3,
			Confidence: math_convert.FromFloat64(0.95),
			Source:     "power",
		},
		internal_environment.Signal{
			Name:       "pci_devices",
			Value:      minFloat(float64(len(fp.PCI))/10.0, 1.0),
			Weight:     0.3,
			Confidence: math_convert.FromFloat64(0.9),
			Source:     "runtime",
		},
		internal_environment.Signal{
			Name:       "mac_addresses",
			Value:      minFloat(float64(len(fp.MAC))/5.0, 1.0),
			Weight:     0.2,
			Confidence: math_convert.FromFloat64(0.9),
			Source:     "runtime",
		},
	)

	// Update Candidates slice for logging
	env.Platform.Candidates = append(env.Platform.Candidates, *s)
}

// ComputeHardwareFingerprint remains unchanged
func ComputeHardwareFingerprint(env *internal_environment.EnvConfig) []byte {
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

	logging.Info("[verification] Environment Hash Generated: %s...", env.Attestation.EnvHash[:12])

	return hash[:]
}

func BoolToFloat(b bool) float64 {
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
