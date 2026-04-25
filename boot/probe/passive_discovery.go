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
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

type SignalBuilder func(fp HardwareFingerprint, env *schema_system.EnvConfig) []schema_system.Signal

type PlatformDefinition struct {
	Type     schema_system.PlatformClass
	Profile  schema_system.PlatformProfile
	Builders []SignalBuilder
}

var platformRegistry = []PlatformDefinition{
	{
		Type: schema_system.PlatformComputer,
		Profile: schema_system.PlatformProfile{
			Class: schema_system.DeviceComputer,
			Form:  schema_system.FormDesktop,
		},
		Builders: []SignalBuilder{buildDesktopSignals},
	},
	{
		Type: schema_system.PlatformMobile,
		Profile: schema_system.PlatformProfile{
			Class: schema_system.DeviceMobile,
			Form:  schema_system.FormPhone,
		},
		Builders: []SignalBuilder{buildMobileSignals},
	},
	{
		Type: schema_system.PlatformEmbedded,
		Profile: schema_system.PlatformProfile{
			Class: schema_system.DeviceEmbedded,
			Form:  schema_system.FormHandheld,
		},
		Builders: []SignalBuilder{buildEmbeddedSignals},
	},
	{
		Type: schema_system.PlatformIndustrial,
		Profile: schema_system.PlatformProfile{
			Class: schema_system.DeviceIndustrial,
		},
		Builders: []SignalBuilder{buildIndustrialSignals},
	},
	{
		Type: schema_system.PlatformVehicle,
		Profile: schema_system.PlatformProfile{
			Class:        schema_system.DeviceVehicle,
			Capabilities: []schema_system.CapabilityTag{schema_system.TagAutomotive},
		},
		Builders: []SignalBuilder{buildVehicleSignals},
	},
	{
		Type: schema_system.PlatformRobot,
		Profile: schema_system.PlatformProfile{
			Class:        schema_system.DeviceRobot,
			Capabilities: []schema_system.CapabilityTag{schema_system.TagDrone},
		},
		Builders: []SignalBuilder{buildRobotSignals},
	},
}

func extractProcessors(fp HardwareFingerprint) []schema_system.Processor {
	var processors []schema_system.Processor

	// CPU
	if n, err := strconv.Atoi(fp.CPU); err == nil && n > 0 {
		processors = append(processors, schema_system.Processor{
			Type:  "CPU",
			Count: runtime.NumCPU(),
		})
	}

	// GPU
	if n, err := strconv.Atoi(fp.GPU); err == nil && n > 0 {
		processors = append(processors, schema_system.Processor{
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

func extractBuses(fp HardwareFingerprint) []schema_system.BusCapability {
	var buses []schema_system.BusCapability

	for k, v := range fp.Buses {
		if !v {
			continue
		}

		buses = append(buses, schema_system.BusCapability{
			ID:         k,
			Type:       strings.ToLower(k),
			Confidence: mathutil.FromFloat64(0.9),
			Source:     "probe",
		})
	}

	return buses
}

func buildHardwareProfile(fp HardwareFingerprint) schema_system.HardwareProfile {
	return schema_system.HardwareProfile{
		Processors: extractProcessors(fp),
		Buses:      extractBuses(fp),
		HasBattery: detectBattery(),
	}
}

// ------------------------------------------------------------
// Public: Passive Discovery
// ------------------------------------------------------------
func PassiveDiscovery(ctx context.Context) (*schema_system.EnvConfig, error) {
	start := time.Now()

	fp, probeErrors := CollectHardwareFingerprint(ctx)
	if len(probeErrors) > 0 {
		logging.Warn("hardware probe errors: %v", probeErrors)
	}

	env := &schema_system.EnvConfig{
		SchemaVersion: 1,
		Identity: schema_system.MachineIdentity{
			OS: runtime.GOOS,
		},
		Hardware:  buildHardwareProfile(fp), // <-- FIX
		Platform:  schema_system.PlatformResolution{},
		Discovery: schema_system.DiscoveryProfile{},
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
func buildRobotSignals(fp HardwareFingerprint, _ *schema_system.EnvConfig) []schema_system.Signal {
	return []schema_system.Signal{
		{Name: "i2c_bus", Value: boolToFloat(hasBus(fp, "i2c")), Weight: 0.5, Confidence: mathutil.FromFloat64(0.9)},
		{Name: "spi_bus", Value: boolToFloat(hasBus(fp, "spi")), Weight: 0.5, Confidence: mathutil.FromFloat64(0.9)},
	}
}

func buildDesktopSignals(fp HardwareFingerprint, env *schema_system.EnvConfig) []schema_system.Signal {
	return []schema_system.Signal{
		{Name: "cpu_cores", Value: minFloat(float64(runtime.NumCPU())/16.0, 1.0), Weight: 0.3, Confidence: mathutil.FromFloat64(0.9)},
		{Name: "pci_devices", Value: minFloat(float64(len(fp.PCI))/10.0, 1.0), Weight: 0.2, Confidence: mathutil.FromFloat64(0.8)},
		{Name: "mac_interfaces", Value: minFloat(float64(len(fp.MAC))/5.0, 1.0), Weight: 0.2, Confidence: mathutil.FromFloat64(0.85)},
		{Name: "battery_present", Value: boolToFloat(env.Hardware.HasBattery), Weight: 0.3, Confidence: mathutil.FromFloat64(0.95)},
	}
}

func buildMobileSignals(fp HardwareFingerprint, env *schema_system.EnvConfig) []schema_system.Signal {
	osName := strings.ToLower(env.Identity.OS)
	return []schema_system.Signal{
		{Name: "mobile_os", Value: boolToFloat(osName == "android" || osName == "ios"), Weight: 0.6, Confidence: mathutil.FromFloat64(0.95)},
		{Name: "battery_present", Value: boolToFloat(env.Hardware.HasBattery), Weight: 0.4, Confidence: mathutil.FromFloat64(0.9)},
	}
}

func buildEmbeddedSignals(fp HardwareFingerprint, _ *schema_system.EnvConfig) []schema_system.Signal {
	return []schema_system.Signal{
		{Name: "low_cpu", Value: 1.0 - minFloat(float64(runtime.NumCPU())/8.0, 1.0), Weight: 0.4, Confidence: mathutil.FromFloat64(0.8)},
		{Name: "has_gpio", Value: boolToFloat(hasBus(fp, "gpio")), Weight: 0.6, Confidence: mathutil.FromFloat64(0.9)},
	}
}

func buildIndustrialSignals(fp HardwareFingerprint, _ *schema_system.EnvConfig) []schema_system.Signal {
	return []schema_system.Signal{
		{Name: "fieldbus", Value: boolToFloat(hasBus(fp, "modbus") || hasBus(fp, "profibus")), Weight: 0.7, Confidence: mathutil.FromFloat64(0.9)},
		{Name: "multi_nic", Value: minFloat(float64(len(fp.MAC))/4.0, 1.0), Weight: 0.3, Confidence: mathutil.FromFloat64(0.8)},
	}
}

func buildVehicleSignals(fp HardwareFingerprint, env *schema_system.EnvConfig) []schema_system.Signal {
	osName := strings.ToLower(env.Identity.OS)
	return []schema_system.Signal{
		{Name: "can_bus", Value: boolToFloat(hasBus(fp, "can")), Weight: 0.6, Confidence: mathutil.FromFloat64(0.95), Source: "bus"},
		{Name: "automotive_os", Value: boolToFloat(osName == "qnx" || osName == "autosar"), Weight: 0.8, Confidence: mathutil.FromFloat64(0.9), Source: "os"},
	}
}

// ------------------------------------------------------------
// PlatformScore
// ------------------------------------------------------------
type LocalPlatformScore struct {
	schema_system.PlatformScore
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

func buildPlatformScore(def PlatformDefinition, fp HardwareFingerprint, env *schema_system.EnvConfig) *LocalPlatformScore {
	var signals []schema_system.Signal
	for _, builder := range def.Builders {
		signals = append(signals, builder(fp, env)...)
	}

	ps := &LocalPlatformScore{
		schema_system.PlatformScore{
			Type:    def.Type,
			Signals: signals,
			Profile: def.Profile,
		},
	}
	ps.Compute()
	return ps
}

func runPlatformInference(env *schema_system.EnvConfig, fp HardwareFingerprint) {
	var results []schema_system.PlatformScore
	var best schema_system.PlatformClass = schema_system.PlatformUnknown
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
	env.Platform.Locked = best != schema_system.PlatformUnknown
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[IDENTITY] Platform: %s (score %.2f)", best, bestScore)
}
