// boot/probe/passive_discovery.go

package probe

import (
	"context"
	"runtime"
	"strconv"
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
	{
		Type: schema.PlatformComputer,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceComputer,
			Form:  schema.FormDesktop,
		},
		Builders: []SignalBuilder{buildDesktopSignals},
	},
	{
		Type: schema.PlatformMobile,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceMobile,
			Form:  schema.FormPhone,
		},
		Builders: []SignalBuilder{buildMobileSignals},
	},
	{
		Type: schema.PlatformEmbedded,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceEmbedded,
			Form:  schema.FormHandheld,
		},
		Builders: []SignalBuilder{buildEmbeddedSignals},
	},
	{
		Type: schema.PlatformIndustrial,
		Profile: schema.PlatformProfile{
			Class: schema.DeviceIndustrial,
		},
		Builders: []SignalBuilder{buildIndustrialSignals},
	},
	{
		Type: schema.PlatformVehicle,
		Profile: schema.PlatformProfile{
			Class:        schema.DeviceVehicle,
			Capabilities: []schema.CapabilityTag{schema.TagAutomotive},
		},
		Builders: []SignalBuilder{buildVehicleSignals},
	},
	{
		Type: schema.PlatformRobot,
		Profile: schema.PlatformProfile{
			Class:        schema.DeviceRobot,
			Capabilities: []schema.CapabilityTag{schema.TagDrone},
		},
		Builders: []SignalBuilder{buildRobotSignals},
	},
}

func extractProcessors(fp HardwareFingerprint) []schema.Processor {
	var processors []schema.Processor

	// CPU
	if n, err := strconv.Atoi(fp.CPU); err == nil && n > 0 {
		processors = append(processors, schema.Processor{
	Type:  "CPU",
	Count: runtime.NumCPU(),
})
	}

	// GPU
	if n, err := strconv.Atoi(fp.GPU); err == nil && n > 0 {
		processors = append(processors, schema.Processor{
			Type:    "GPU",
			Count:   n,
			Version: 1.0,
		})
	}

	return processors
}

func extractBuses(fp HardwareFingerprint) []schema.BusCapability {
	var buses []schema.BusCapability

	for k, v := range fp.Buses {
		if !v {
			continue
		}

		buses = append(buses, schema.BusCapability{
			ID:         k,
			Type:       strings.ToLower(k),
			Confidence: mathutil.FromFloat64(0.9),
			Source:     "probe",
		})
	}

	return buses
}

func buildHardwareProfile(fp HardwareFingerprint) schema.HardwareProfile {
	return schema.HardwareProfile{
		Processors: extractProcessors(fp),
		Buses:      extractBuses(fp),
		HasBattery: detectBattery(),
	}
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
		SchemaVersion: 1,
		Identity: schema.MachineIdentity{
			OS: runtime.GOOS,
		},
		Hardware:  buildHardwareProfile(fp), // <-- FIX
		Platform:  schema.PlatformResolution{},
		Discovery: schema.DiscoveryProfile{},
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
func buildRobotSignals(fp HardwareFingerprint, _ *schema.EnvConfig) []schema.Signal {
	return []schema.Signal{
		{Name: "i2c_bus", Value: boolToFloat(hasBus(fp, "i2c")), Weight: 0.5, Confidence: mathutil.FromFloat64(0.9)},
		{Name: "spi_bus", Value: boolToFloat(hasBus(fp, "spi")), Weight: 0.5, Confidence: mathutil.FromFloat64(0.9)},
	}
}

func buildDesktopSignals(fp HardwareFingerprint, env *schema.EnvConfig) []schema.Signal {
	return []schema.Signal{
		{Name: "cpu_cores", Value: minFloat(float64(runtime.NumCPU())/16.0, 1.0), Weight: 0.3, Confidence: mathutil.FromFloat64(0.9)},
		{Name: "pci_devices", Value: minFloat(float64(len(fp.PCI))/10.0, 1.0), Weight: 0.2, Confidence: mathutil.FromFloat64(0.8)},
		{Name: "mac_interfaces", Value: minFloat(float64(len(fp.MAC))/5.0, 1.0), Weight: 0.2, Confidence: mathutil.FromFloat64(0.85)},
		{Name: "battery_present", Value: boolToFloat(env.Hardware.HasBattery), Weight: 0.3, Confidence: mathutil.FromFloat64(0.95)},
	}
}

func buildMobileSignals(fp HardwareFingerprint, env *schema.EnvConfig) []schema.Signal {
	osName := strings.ToLower(env.Identity.OS)
	return []schema.Signal{
		{Name: "mobile_os", Value: boolToFloat(osName == "android" || osName == "ios"), Weight: 0.6, Confidence: mathutil.FromFloat64(0.95)},
		{Name: "battery_present", Value: boolToFloat(env.Hardware.HasBattery), Weight: 0.4, Confidence: mathutil.FromFloat64(0.9)},
	}
}

func buildEmbeddedSignals(fp HardwareFingerprint, _ *schema.EnvConfig) []schema.Signal {
	return []schema.Signal{
		{Name: "low_cpu", Value: 1.0 - minFloat(float64(runtime.NumCPU())/8.0, 1.0), Weight: 0.4, Confidence: mathutil.FromFloat64(0.8)},
		{Name: "has_gpio", Value: boolToFloat(hasBus(fp, "gpio")), Weight: 0.6, Confidence: mathutil.FromFloat64(0.9)},
	}
}

func buildIndustrialSignals(fp HardwareFingerprint, _ *schema.EnvConfig) []schema.Signal {
	return []schema.Signal{
		{Name: "fieldbus", Value: boolToFloat(hasBus(fp, "modbus") || hasBus(fp, "profibus")), Weight: 0.7, Confidence: mathutil.FromFloat64(0.9)},
		{Name: "multi_nic", Value: minFloat(float64(len(fp.MAC))/4.0, 1.0), Weight: 0.3, Confidence: mathutil.FromFloat64(0.8)},
	}
}

func buildVehicleSignals(fp HardwareFingerprint, env *schema.EnvConfig) []schema.Signal {
	osName := strings.ToLower(env.Identity.OS)
	return []schema.Signal{
		{Name: "can_bus", Value: boolToFloat(hasBus(fp, "can")), Weight: 0.6, Confidence: mathutil.FromFloat64(0.95), Source: "bus"},
		{Name: "automotive_os", Value: boolToFloat(osName == "qnx" || osName == "autosar"), Weight: 0.8, Confidence: mathutil.FromFloat64(0.9), Source: "os"},
	}
}

// ------------------------------------------------------------
// PlatformScore
// ------------------------------------------------------------
type LocalPlatformScore struct {
	schema.PlatformScore
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
		ps.Confidence = mathutil.FromFloat64(total / max)
	} else {
		ps.Confidence = mathutil.Min
	}

	for _, s := range ps.Signals {
		logging.Info("[SIGNAL] %s val=%.2f weight=%.2f conf=%.2f",
			s.Name, s.Value, s.Weight, s.Confidence.Float64())
	}
}

func buildPlatformScore(def PlatformDefinition, fp HardwareFingerprint, env *schema.EnvConfig) *LocalPlatformScore {
	var signals []schema.Signal
	for _, builder := range def.Builders {
		signals = append(signals, builder(fp, env)...)
	}

	ps := &LocalPlatformScore{
		schema.PlatformScore{
			Type:    def.Type,
			Signals: signals,
			Profile: def.Profile,
		},
	}
	ps.Compute()
	return ps
}

func runPlatformInference(env *schema.EnvConfig, fp HardwareFingerprint) {
	var results []schema.PlatformScore
	var best schema.PlatformClass = schema.PlatformUnknown
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
	env.Platform.Locked = best != schema.PlatformUnknown
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[IDENTITY] Platform: %s (score %.2f)", best, bestScore)
}
