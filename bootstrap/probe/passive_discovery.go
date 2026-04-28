// bootstrap/probe/passive_discovery.go

package probe

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/math_convert"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

type SignalBuilder func(fp HardwareFingerprint, env *internal_environment.EnvConfig) []internal_environment.Signal

type PlatformDefinition struct {
	Type     internal_environment.PlatformClass
	Profile  internal_environment.PlatformProfile
	Builders []SignalBuilder
}

var platformRegistry = []PlatformDefinition{
	{
		Type: internal_environment.PlatformComputer,
		Profile: internal_environment.PlatformProfile{
			Class: internal_environment.DeviceComputer,
			Form:  internal_environment.FormDesktop,
		},
		Builders: []SignalBuilder{buildDesktopSignals},
	},
	{
		Type: internal_environment.PlatformMobile,
		Profile: internal_environment.PlatformProfile{
			Class: internal_environment.DeviceMobile,
			Form:  internal_environment.FormPhone,
		},
		Builders: []SignalBuilder{buildMobileSignals},
	},
	{
		Type: internal_environment.PlatformEmbedded,
		Profile: internal_environment.PlatformProfile{
			Class: internal_environment.DeviceEmbedded,
			Form:  internal_environment.FormHandheld,
		},
		Builders: []SignalBuilder{buildEmbeddedSignals},
	},
	{
		Type: internal_environment.PlatformIndustrial,
		Profile: internal_environment.PlatformProfile{
			Class: internal_environment.DeviceIndustrial,
		},
		Builders: []SignalBuilder{buildIndustrialSignals},
	},
	{
		Type: internal_environment.PlatformVehicle,
		Profile: internal_environment.PlatformProfile{
			Class:        internal_environment.DeviceVehicle,
			Capabilities: []internal_environment.CapabilityTag{internal_environment.TagAutomotive},
		},
		Builders: []SignalBuilder{buildVehicleSignals},
	},
	{
		Type: internal_environment.PlatformRobot,
		Profile: internal_environment.PlatformProfile{
			Class:        internal_environment.DeviceRobot,
			Capabilities: []internal_environment.CapabilityTag{internal_environment.TagDrone},
		},
		Builders: []SignalBuilder{buildRobotSignals},
	},
}

func extractProcessors(fp HardwareFingerprint) []internal_environment.Processor {
	var processors []internal_environment.Processor

	// CPU
	if n, err := strconv.Atoi(fp.CPU); err == nil && n > 0 {
		processors = append(processors, internal_environment.Processor{
			Type:  "CPU",
			Count: runtime.NumCPU(),
		})
	}

	// GPU
	if n, err := strconv.Atoi(fp.GPU); err == nil && n > 0 {
		processors = append(processors, internal_environment.Processor{
			Type:    "GPU",
			Count:   NumGPU(),
			Version: 1.0,
		})
	}

	return processors
}
func NumGPU() int {
	return 0 // placeholder or platform-specific detection
}

func extractBuses(fp HardwareFingerprint) []internal_environment.BusCapability {
	var buses []internal_environment.BusCapability

	for k, v := range fp.Buses {
		if !v {
			continue
		}

		buses = append(buses, internal_environment.BusCapability{
			ID:         k,
			Type:       strings.ToLower(k),
			Confidence: math_convert.FromFloat64(0.9),
			Source:     "probe",
		})
	}

	return buses
}

func buildHardwareProfile(fp HardwareFingerprint) internal_environment.HardwareProfile {
	return internal_environment.HardwareProfile{
		Processors: extractProcessors(fp),
		Buses:      extractBuses(fp),
		HasBattery: detectBattery(),
	}
}

// ------------------------------------------------------------
// Public: Passive Discovery
// ------------------------------------------------------------
func PassiveDiscovery(ctx context.Context) (*internal_environment.EnvConfig, error) {
	start := time.Now()

	fp, probeErrors := CollectHardwareFingerprint(ctx)
	if len(probeErrors) > 0 {
		logging.Warn("hardware probe errors: %v", probeErrors)
	}

	env := &internal_environment.EnvConfig{
		SchemaVersion: 1,
		Identity: internal_environment.MachineIdentity{
			OS: runtime.GOOS,
		},
		Hardware:  buildHardwareProfile(fp), // <-- FIX
		Platform:  internal_environment.PlatformResolution{},
		Discovery: internal_environment.DiscoveryProfile{},
	}

	runPlatformInference(env, fp)

	duration := time.Since(start)
	env.Discovery.DiscoveryDuration = duration
	logging.Info("[DISCOVERY] Duration: %s", duration)

	return env, nil
}

// ------------------------------------------------------------
// Helpers
// ------------------------------------------------------------
func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

func hasBus(fp HardwareFingerprint, busType string) bool {
	// Use case-insensitive key matching
	for k, v := range fp.Buses {
		if strings.EqualFold(k, busType) && v {
			return true
		}
	}
	return false
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// ------------------------------------------------------------
// Signal Builders
// ------------------------------------------------------------
func buildRobotSignals(fp HardwareFingerprint, _ *internal_environment.EnvConfig) []internal_environment.Signal {
	return []internal_environment.Signal{
		{Name: "i2c_bus", Value: boolToFloat(hasBus(fp, "i2c")), Weight: 0.5, Confidence: math_convert.FromFloat64(0.9)},
		{Name: "spi_bus", Value: boolToFloat(hasBus(fp, "spi")), Weight: 0.5, Confidence: math_convert.FromFloat64(0.9)},
	}
}

func buildDesktopSignals(fp HardwareFingerprint, env *internal_environment.EnvConfig) []internal_environment.Signal {
	return []internal_environment.Signal{
		{Name: "cpu_cores", Value: minFloat(float64(runtime.NumCPU())/16.0, 1.0), Weight: 0.3, Confidence: math_convert.FromFloat64(0.9)},
		{Name: "pci_devices", Value: minFloat(float64(len(fp.PCI))/10.0, 1.0), Weight: 0.2, Confidence: math_convert.FromFloat64(0.8)},
		{Name: "mac_interfaces", Value: minFloat(float64(len(fp.MAC))/5.0, 1.0), Weight: 0.2, Confidence: math_convert.FromFloat64(0.85)},
		{Name: "battery_present", Value: boolToFloat(env.Hardware.HasBattery), Weight: 0.3, Confidence: math_convert.FromFloat64(0.95)},
	}
}

func buildMobileSignals(fp HardwareFingerprint, env *internal_environment.EnvConfig) []internal_environment.Signal {
	osName := strings.ToLower(env.Identity.OS)
	return []internal_environment.Signal{
		{Name: "mobile_os", Value: boolToFloat(osName == "android" || osName == "ios"), Weight: 0.6, Confidence: math_convert.FromFloat64(0.95)},
		{Name: "battery_present", Value: boolToFloat(env.Hardware.HasBattery), Weight: 0.4, Confidence: math_convert.FromFloat64(0.9)},
	}
}

func buildEmbeddedSignals(fp HardwareFingerprint, _ *internal_environment.EnvConfig) []internal_environment.Signal {
	return []internal_environment.Signal{
		{Name: "low_cpu", Value: 1.0 - minFloat(float64(runtime.NumCPU())/8.0, 1.0), Weight: 0.4, Confidence: math_convert.FromFloat64(0.8)},
		{Name: "has_gpio", Value: boolToFloat(hasBus(fp, "gpio")), Weight: 0.6, Confidence: math_convert.FromFloat64(0.9)},
	}
}

func buildIndustrialSignals(fp HardwareFingerprint, _ *internal_environment.EnvConfig) []internal_environment.Signal {
	return []internal_environment.Signal{
		{Name: "fieldbus", Value: boolToFloat(hasBus(fp, "modbus") || hasBus(fp, "profibus")), Weight: 0.7, Confidence: math_convert.FromFloat64(0.9)},
		{Name: "multi_nic", Value: minFloat(float64(len(fp.MAC))/4.0, 1.0), Weight: 0.3, Confidence: math_convert.FromFloat64(0.8)},
	}
}

func buildVehicleSignals(fp HardwareFingerprint, env *internal_environment.EnvConfig) []internal_environment.Signal {
	osName := strings.ToLower(env.Identity.OS)
	return []internal_environment.Signal{
		{Name: "can_bus", Value: boolToFloat(hasBus(fp, "can")), Weight: 0.6, Confidence: math_convert.FromFloat64(0.95), Source: "bus"},
		{Name: "automotive_os", Value: boolToFloat(osName == "qnx" || osName == "autosar"), Weight: 0.8, Confidence: math_convert.FromFloat64(0.9), Source: "os"},
	}
}

// ------------------------------------------------------------
// PlatformScore
// ------------------------------------------------------------
type LocalPlatformScore struct {
	internal_environment.PlatformScore
}

func (ps *LocalPlatformScore) Compute() {
	var total, max float64
	for _, s := range ps.Signals {
		total += s.Value * s.Weight * s.Confidence.Float64()
		max += s.Weight
	}

	ps.Score = total
	ps.MaxScore = max
	if max > 0 {
		ps.Confidence = math_convert.FromFloat64(total / max)
	} else {
		ps.Confidence = math_convert.Min
	}

	for _, s := range ps.Signals {
		logging.Info("[SIGNAL] %s val=%.2f weight=%.2f conf=%.2f",
			s.Name, s.Value, s.Weight, s.Confidence.Float64())
	}
}

func buildPlatformScore(def PlatformDefinition, fp HardwareFingerprint, env *internal_environment.EnvConfig) *LocalPlatformScore {
	var signals []internal_environment.Signal
	for _, builder := range def.Builders {
		signals = append(signals, builder(fp, env)...)
	}

	ps := &LocalPlatformScore{
		internal_environment.PlatformScore{
			Type:    def.Type,
			Signals: signals,
			Profile: def.Profile,
		},
	}
	ps.Compute()
	return ps
}

func runPlatformInference(env *internal_environment.EnvConfig, fp HardwareFingerprint) {
	var results []internal_environment.PlatformScore
	var best internal_environment.PlatformClass = internal_environment.PlatformUnknown
	var bestScore float64
	const minConfidence = 0.65
	const delta = 0.1 // margin between top candidates

	for _, def := range platformRegistry {
		ps := buildPlatformScore(def, fp, env)
		results = append(results, ps.PlatformScore)

		conf := ps.Confidence.Float64()
		if conf >= minConfidence && conf > bestScore {
			bestScore = conf
			best = ps.Type
		}
	}

	env.Platform.Candidates = results
	env.Platform.Final = best
	env.Platform.Locked = best != internal_environment.PlatformUnknown
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[IDENTITY] Platform: %s (score %.2f)", best, bestScore)
}
