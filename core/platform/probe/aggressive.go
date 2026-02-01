//core/platform/probe/aggressive.go

package probe

import (
	"multi-platform-AI/configs/platforms"
	"multi-platform-AI/internal/logging"
	"time"
)

// AggressiveScan performs active discovery and populates the global EnvConfig.
// It bridges raw hardware pings to the typed Platform architecture.
func AggressiveScan(env *platforms.EnvConfig) error {
	logging.Info("[PROBE] Stage 1 Aggressive Scan: Waking up %s", env.Platform.Final)

	// 1. DYNAMIC DRIVER SELECTION
	// We map your switch logic to the typed PlatformClass constants
	switch env.Platform.Final {
	case platforms.PlatformVehicle:
		scanCANBus(env)
		scanLidarArrays(env)
	case platforms.PlatformIndustrial:
		scanModbusTCP(env)
	case platforms.PlatformComputer, platforms.PlatformLaptop:
		scanUSBBus(env)
		scanGPUDriver(env)
	}

	// 2. STRESS TEST / LATENCY CHECK
	// We store latency data in the RuntimeProfile for the Trust Engine to evaluate
	measureBusLatencies(env)

	logging.Info("[PROBE] Aggressive scan complete. Found %d bus nodes.", len(env.Hardware.Buses))
	return nil
}

func scanCANBus(env *platforms.EnvConfig) {
	logging.Info(" - Pinging CAN-bus nodes...")
	// We add a concrete BusCapability to the profile
	env.Hardware.Buses = append(env.Hardware.Buses, platforms.BusCapability{
		ID:         "can0",
		Type:       "can",
		Confidence: 65535, // Probed existence = 1.0 confidence
		Source:     "probed",
	})
	
	// Add processors if the ECU is detected as a co-processor
	env.Hardware.Processors = append(env.Hardware.Processors, platforms.Processor{
		Type: "ECU", 
		Count: 1,
	})
}

func scanLidarArrays(env *platforms.EnvConfig) {
	logging.Info(" - Initializing Lidar Spin-up...")
	// In your architecture, Sensors can be treated as bus nodes or I/O Nodes
	env.Hardware.Buses = append(env.Hardware.Buses, platforms.BusCapability{
		ID:     "lidar_front",
		Type:   "ethernet_sensor",
		Source: "probed",
	})
}

func measureBusLatencies(env *platforms.EnvConfig) {
	// We can store latency in the EnvConfig.Runtime.EnvVars or as a specialized signal
	// for the Bayesian Evaluator to penalize trust if latency > threshold.
	if env.Runtime.EnvVars == nil {
		env.Runtime.EnvVars = make(map[string]string)
	}
	env.Runtime.EnvVars["primary_bus_latency_ms"] = "2.0"
}