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

type SignalBuilder func(fp HardwareFingerprint, env *schema.EnvConfig) []schema.Signal

type PlatformDefinition struct {
	Type     schema.PlatformClass
	Profile  schema.PlatformProfile
	Builders []SignalBuilder
}


var platformRegistry = []PlatformDefinition{

	// --------------------------------------------------
	// COMPUTER
	// --------------------------------------------------
	{
		Type: schema.PlatformComputer,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceComputer,
			Form:  schema.FormDesktop,
		},
		Builders: []SignalBuilder{
			buildDesktopSignals,
		},
	},

	// --------------------------------------------------
	// MOBILE
	// --------------------------------------------------
	{
		Type: schema.PlatformMobile,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceMobile,
			Form:  schema.FormPhone,
		},
		Builders: []SignalBuilder{
			buildMobileSignals,
		},
	},

	// --------------------------------------------------
	// EMBEDDED
	// --------------------------------------------------
	{
		Type: schema.PlatformEmbedded,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceEmbedded,
			Form:  schema.FormHandheld,
		},
		Builders: []SignalBuilder{
			buildEmbeddedSignals,
		},
	},

	// --------------------------------------------------
	// INDUSTRIAL
	// --------------------------------------------------
	{
		Type: schema.PlatformIndustrial,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceIndustrial,
		},
		Builders: []SignalBuilder{
			buildIndustrialSignals,
		},
	},

	// --------------------------------------------------
	// VEHICLE
	// --------------------------------------------------
	{
		Type: schema.PlatformVehicle,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceVehicle,
			Capabilities: []schema.CapabilityTag{
				schema.TagAutomotive,
			},
		},
		Builders: []SignalBuilder{
			buildVehicleSignals,
		},
	},

	// --------------------------------------------------
	// ROBOT
	// --------------------------------------------------
	{
		Type: schema.PlatformRobot,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceRobot,
			Capabilities: []schema.CapabilityTag{
				schema.TagDrone, // optional overlap
			},
		},
		Builders: []SignalBuilder{
			buildRobotSignals,
		},
	},
}

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

func buildVehicleSignals(fp HardwareFingerprint, env *schema.EnvConfig) []schema.Signal {
	osName := strings.ToLower(env.Identity.OS)

	return []schema.Signal{
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
}

func buildRobotSignals(fp HardwareFingerprint, _ *schema.EnvConfig) []schema.Signal {
	return []schema.Signal{
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
}

func buildDesktopSignals(fp HardwareFingerprint, env *schema.EnvConfig) []schema.Signal {
	return []schema.Signal{
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
}

func buildPlatformScore(def PlatformDefinition, fp HardwareFingerprint, env *schema.EnvConfig) *schema.PlatformScore {
	var signals []schema.Signal

	for _, b := range def.Builders {
		signals = append(signals, b(fp, env)...)
	}

	return &schema.PlatformScore{
		Type:    def.Type,
		Signals: signals,
	}
}

func buildMobileSignals(fp HardwareFingerprint, env *schema.EnvConfig) []schema.Signal {
	os := strings.ToLower(env.Identity.OS)

	return []schema.Signal{
		{
			Name:       "mobile_os",
			Value:      boolToFloat(os == "android" || os == "ios"),
			Weight:     0.6,
			Confidence: 0.95,
		},
		{
			Name:       "battery_present",
			Value:      boolToFloat(env.Hardware.HasBattery),
			Weight:     0.4,
			Confidence: 0.9,
		},
	}
}

func buildEmbeddedSignals(fp HardwareFingerprint, env *schema.EnvConfig) []schema.Signal {
	return []schema.Signal{
		{
			Name:       "low_cpu",
			Value:      1.0 - minFloat(float64(runtime.NumCPU())/8.0, 1.0),
			Weight:     0.4,
			Confidence: 0.8,
		},
		{
			Name:       "has_gpio",
			Value:      boolToFloat(hasBus(fp, "gpio")),
			Weight:     0.6,
			Confidence: 0.9,
		},
	}
}

func buildIndustrialSignals(fp HardwareFingerprint, _ *schema.EnvConfig) []schema.Signal {
	return []schema.Signal{
		{
			Name:       "fieldbus",
			Value:      boolToFloat(hasBus(fp, "modbus") || hasBus(fp, "profibus")),
			Weight:     0.7,
			Confidence: 0.9,
		},
		{
			Name:       "multi_nic",
			Value:      minFloat(float64(len(fp.MAC))/4.0, 1.0),
			Weight:     0.3,
			Confidence: 0.8,
		},
	}
}

func runPlatformInference(env *schema.EnvConfig, fp HardwareFingerprint) {
	var results []schema.PlatformScore

	var best schema.PlatformClass = schema.PlatformUnknown
	var bestScore float64

	for _, def := range platformRegistry {
		ps := buildPlatformScore(def, fp, env)
		ps.Profile = def.Profile
		ps.Compute()
		ps.Q16 = mathutil.FromFloat64(ps.Confidence)

		results = append(results, *ps)

		if ps.Confidence > bestScore && ps.Confidence >= minConfidence {
			bestScore = ps.Confidence
			best = ps.Type
		}
	}

	env.Platform.Candidates = results
	env.Platform.Final = best
	env.Platform.Locked = best != schema.PlatformUnknown
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[IDENTITY] Platform: %s (%.2f)", best, bestScore)
}