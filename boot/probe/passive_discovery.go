// boot/probe/passive_discovery.go
package probe

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ------------------------------------------------------------
// Public: Passive Discovery
// ------------------------------------------------------------

func PassiveDiscovery(ctx context.Context) (*schema.EnvConfig, error) {
	start := time.Now()

	fp, probeErrors := CollectHardwareFingerprint(ctx)
	if len(probeErrors) > 0 {
		logging.Warn("hardware probe errors: %v", probeErrors)
	}

	env := &schema.EnvConfig{
		Identity: schema.Identity{
			OS: runtime.GOOS,
		},
	}

	runPlatformInference(env, fp)

	if env.Discovery == nil {
		env.Discovery = &schema.DiscoveryDiagnostics{}
	}
	env.Discovery.DiscoveryDuration = time.Since(start)

	return env, nil
}

func runProbe[T any](ctx context.Context, name string, fn func(context.Context) (T, error)) ProbeResult[T] {
	start := time.Now()
	val, err := fn(ctx)

	return ProbeResult[T]{
		Value:    val,
		Error:    err,
		Duration: time.Since(start),
		Source:   name,
	}
}

//
// ------------------------------------------------------------
// Machine Identity
// ------------------------------------------------------------
//

func BuildRobustMachineID(fp HardwareFingerprint) string {
	stable := []string{fp.TPM, fp.DMI}
	semi := []string{fp.CPU, strings.Join(fp.Storage, ",")}
	volatile := strings.Join(fp.MAC, "|")

	hasAnchor := false
	for _, s := range stable {
		if strings.TrimSpace(s) != "" {
			hasAnchor = true
			break
		}
	}

	var material string

	switch {
	case hasAnchor:
		material = strings.Join(stable, "|") + "::" + strings.Join(semi, "|")

	case volatile != "":
		material = volatile + "::" + strings.Join(semi, "|")

	default:
		hostname, _ := os.Hostname()
		material = hostname + "::" + runtime.GOOS
	}

	hash := sha256.Sum256([]byte(material))
	return hex.EncodeToString(hash[:])
}
//
// ------------------------------------------------------------
// Platform Inference
// ------------------------------------------------------------
//

const minConfidence = 0.25

func runPlatformInference(env *schema.EnvConfig, fp HardwareFingerprint) {
	var candidates []*schema.PlatformScore

	osName := strings.ToLower(env.Identity.OS)

	candidates = append(candidates,
		scoreVehicle(fp, osName),
		scoreRobot(fp),
		scoreDesktop(fp, env),
	)

	var best schema.PlatformClass = schema.PlatformUnknown
	var bestScore float64

	var results []schema.PlatformScore

	for _, c := range candidates {
		if c == nil {
			continue
		}

		c.Compute()

		c.ConfidenceQ16 = mathutil.FromFloat64(c.Confidence)

		results = append(results, *c)

		if c.Confidence > bestScore && c.Confidence >= minConfidence {
			bestScore = c.Confidence
			best = c.Type
		}
	}

	env.Platform.Candidates = results
	env.Platform.Final = best
	env.Platform.Locked = best != schema.PlatformUnknown
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[IDENTITY] Platform: %s (%.2f)", best, bestScore)
}

//
// ------------------------------------------------------------
// Desktop Scoring
// ------------------------------------------------------------
//

func collectDesktopSignals(fp HardwareFingerprint, env *schema.EnvConfig) *schema.PlatformScore {

	signals := []schema.Signal{
		{
			Name:       "cpu_cores",
			Value:      minFloat(float64(runtime.NumCPU())/16.0, 1.0),
			Weight:     0.3,
			Confidence: 0.9,
			Source:     "runtime",
		},
		{
			Name:       "pci_devices",
			Value:      minFloat(float64(len(fp.PCI))/10.0, 1.0),
			Weight:     0.2,
			Confidence: 0.8,
			Source:     "lspci",
		},
		{
			Name:       "mac_interfaces",
			Value:      minFloat(float64(len(fp.MAC))/5.0, 1.0),
			Weight:     0.2,
			Confidence: 0.85,
			Source:     "net",
		},
		{
			Name:       "battery_present",
			Value:      boolToFloat(env.Hardware.HasBattery),
			Weight:     0.3,
			Confidence: 0.95,
			Source:     "power",
		},
	}

	ps := &schema.PlatformScore{
		Type:    schema.PlatformComputer,
		Signals: signals,
	}

	ps.Compute()
	return ps
}

//
// ------------------------------------------------------------
// Helpers
// ------------------------------------------------------------
//

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (ps *schema.PlatformScore) Compute() {
	var total float64
	var max float64

	for _, s := range ps.Signals {
		contrib := s.Value * s.Weight * s.Confidence
		total += contrib
		max += s.Weight
	}

	ps.Score = total
	ps.MaxScore = max

	if max > 0 {
		ps.Confidence = total / max
	}
}

func vehicleSignals(fp HardwareFingerprint, osName string) []schema.Signal {
	return []schema.Signal{
		{
			Name:       "can_bus",
			Value:      boolToFloat(hasBus(fp, "can")),
			Weight:     0.6,
			Confidence: 0.95,
			Source:     "bus-registry",
		},
		{
			Name:       "automotive_os",
			Value:      boolToFloat(osName == "qnx" || osName == "autosar"),
			Weight:     0.8,
			Confidence: 0.9,
			Source:     "os",
		},
	}
}


func scoreVehicle(fp HardwareFingerprint, osName string) *schema.PlatformScore {
	signals := []schema.Signal{
		{
			Name:       "can_bus",
			Value:      boolToFloat(hasBus(fp, "can")),
			Weight:     0.6,
			Confidence: 0.95,
			Source:     "bus",
		},
		{
			Name:       "automotive_os",
			Value:      boolToFloat(osName == "qnx" || osName == "autosar"),
			Weight:     0.8,
			Confidence: 0.9,
			Source:     "os",
		},
	}

	return &schema.PlatformScore{
		Type:    schema.PlatformVehicle,
		Signals: signals,
	}
}

func scoreRobot(fp HardwareFingerprint) *schema.PlatformScore {
	signals := []schema.Signal{
		{
			Name:       "i2c_bus",
			Value:      boolToFloat(hasBus(fp, "i2c")),
			Weight:     0.5,
			Confidence: 0.9,
		},
		{
			Name:       "spi_bus",
			Value:      boolToFloat(hasBus(fp, "spi")),
			Weight:     0.5,
			Confidence: 0.9,
		},
	}

	return &schema.PlatformScore{
		Type:    schema.PlatformRobot,
		Signals: signals,
	}
}

func scoreDesktop(fp HardwareFingerprint, env *schema.EnvConfig) *schema.PlatformScore {
	signals := []schema.Signal{
		{
			Name:       "cpu_cores",
			Value:      minFloat(float64(runtime.NumCPU())/16.0, 1.0),
			Weight:     0.3,
			Confidence: 0.9,
		},
		{
			Name:       "pci_devices",
			Value:      minFloat(float64(len(fp.PCI))/10.0, 1.0),
			Weight:     0.2,
			Confidence: 0.8,
		},
		{
			Name:       "mac_interfaces",
			Value:      minFloat(float64(len(fp.MAC))/5.0, 1.0),
			Weight:     0.2,
			Confidence: 0.85,
		},
		{
			Name:       "battery_present",
			Value:      boolToFloat(env.Hardware.HasBattery),
			Weight:     0.3,
			Confidence: 0.95,
		},
	}

	return &schema.PlatformScore{
		Type:    schema.PlatformComputer,
		Signals: signals,
	}
}