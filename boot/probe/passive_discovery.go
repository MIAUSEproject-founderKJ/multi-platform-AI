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

env := schema.EnvConfig{

	runPlatformInference(env, fp)

	// Safe diagnostics assignment
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
		// last-resort entropy
		material = fmt.Sprintf("fallback-%d-%s", time.Now().UnixNano(), runtime.GOOS)
	}

	hash := sha256.Sum256([]byte(material))
	return hex.EncodeToString(hash[:])
}

//
// ------------------------------------------------------------
// Platform Inference
// ------------------------------------------------------------
//

func runPlatformInference(env *schema.EnvConfig, fp HardwareFingerprint) {

	scores := map[schema.PlatformClass]*schema.PlatformScore{}

	ensure := func(class schema.PlatformClass, max float64) *schema.PlatformScore {
		if scores[class] == nil {
			scores[class] = &schema.PlatformScore{
				Type:     class,
				MaxScore: max,
			}
		}
		return scores[class]
	}

	osName := strings.ToLower(env.Identity.OS)

	// Vehicle
	if hasBus(fp, "can") || osName == "qnx" || osName == "autosar" {
		s := ensure(schema.PlatformVehicle, 1.5)
		s.Score += 1.0
		s.Signals = append(s.Signals, "automotive environment detected")
	}

	// Robot
	if hasBus(fp, "i2c") && hasBus(fp, "spi") {
		s := ensure(schema.PlatformRobot, 1.2)
		s.Score += 0.4
		s.Signals = append(s.Signals, "sensor buses detected")
	}

	// Desktop
	scores[schema.PlatformComputer] = collectDesktopSignals(fp, env)

	// Resolve best
	var best schema.PlatformClass = schema.PlatformUnknown
	highConf := mathutil.Q16(0)

	var candidates []schema.PlatformScore

	if len(scores) == 0 {
		logging.Warn("[IDENTITY] No platform signals detected")
		env.Platform.Final = schema.PlatformUnknown
		env.Platform.Locked = false
		return
	}

	for _, s := range scores {
		s.Compute()

		s.Confidence = mathutil.Q16(mathutil.FromFloat64(s.Confidence))

		candidates = append(candidates, *s)

		if s.Confidence > highConf {
			highConf = s.Confidence
			best = s.Type
		}
	}

	env.Platform.Candidates = candidates
	env.Platform.Final = best
	env.Platform.Locked = true
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[IDENTITY] Platform: %s (%d%%)", best, highConf.Percentage())
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
		Type:    schema.PlatformVehicle,
		Signals: vehicleSignals(fp, osName),
	}
	ps.Compute()
	scores[schema.PlatformVehicle] = ps
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
