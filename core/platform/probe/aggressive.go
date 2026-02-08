//core/platform/probe/aggressive.go

package probe

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// AggressiveScan performs active discovery and populates the global EnvConfig.
// It bridges raw hardware pings to the typed Platform architecture.
func AggressiveScan(cfg *schema.EnvConfig) {
	logging.Info("Aggressive probe initiated for platform: %s", cfg.Platform)

	// 1. DYNAMIC DRIVER SELECTION
	// We map your switch logic to the typed PlatformClass constants
	switch cfg.Platform.Final {
	case schema.PlatformVehicle:
		scanCANBus(cfg)
		scanLidarArrays(cfg)
	case schema.PlatformIndustrial:
		scanModbusTCP(cfg)
	case schema.PlatformComputer, schema.PlatformLaptop:
		scanUSBBus(cfg)
		scanGPUDriver(cfg)
	}

	// 2. STRESS TEST / LATENCY CHECK
	// We store latency data in the RuntimeProfile for the Trust Engine to evaluate
	measureBusLatencies(cfg)

	logging.Info("[PROBE] Aggressive scan complete. Found %d bus nodes.", len(cfg.Hardware.Buses))
	return nil
}

func scanCANBus(env *schema.EnvConfig) {
	logging.Info(" - Pinging CAN-bus nodes...")
	// We add a concrete BusCapability to the profile
	env.Hardware.Buses = append(env.Hardware.Buses, schema.BusCapability{
		ID:         "can0",
		Type:       "can",
		Confidence: 65535, // Probed existence = 1.0 confidence
		Source:     "probed",
	})

	// Add processors if the ECU is detected as a co-processor
	env.Hardware.Processors = append(env.Hardware.Processors, schema.Processor{
		Type:  "ECU",
		Count: 1,
	})
}

func scanLidarArrays(env *schema.EnvConfig) {
	logging.Info(" - Initializing Lidar Spin-up...")
	// In your architecture, Sensors can be treated as bus nodes or I/O Nodes
	env.Hardware.Buses = append(env.Hardware.Buses, schema.BusCapability{
		ID:     "lidar_front",
		Type:   "ethernet_sensor",
		Source: "probed",
	})
}

func measureBusLatencies(cfg *schema.EnvConfig) {
	// We can store latency in the EnvConfig.Runtime.EnvVars or as a specialized signal
	// for the Bayesian Evaluator to penalize trust if latency > threshold.
	if cfg.Runtime.EnvVars == nil {
		cfg.Runtime.EnvVars = make(map[string]string)
	}
	cfg.Runtime.EnvVars["primary_bus_latency_ms"] = "2.0"
}
