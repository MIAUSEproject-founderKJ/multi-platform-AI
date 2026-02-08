//core/sim/simulation_engine.go

package sim

import (
	"math/rand"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/configs/platforms"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type Scenario string

const (
	ScenarioNormal     Scenario = "nominal_operations"
	ScenarioLidarFail  Scenario = "sensor_blackout"
	ScenarioBusFlood   Scenario = "can_bus_congestion"
	ScenarioIdentityHijack Scenario = "platform_mismatch"
)

type SimulationEngine struct {
	ActiveScenario Scenario
	Intensity      float64 // 0.0 to 1.0 (How severe the fault is)
}

// InjectFault modifies the EnvConfig to simulate real-world hardware issues.
func (se *SimulationEngine) InjectFault(env *platforms.EnvConfig) {
	logging.Warn("[SIM] Injecting Scenario: %s (Intensity: %.2f)", se.ActiveScenario, se.Intensity)

	switch se.ActiveScenario {
	case ScenarioLidarFail:
		// Remove Lidar from detected buses to see if Trust drops
		for i, b := range env.Hardware.Buses {
			if b.Type == "ethernet_sensor" {
				env.Hardware.Buses = append(env.Hardware.Buses[:i], env.Hardware.Buses[i+1:]...)
				logging.Info("[SIM] Fault: Lidar Connection Severed")
				break
			}
		}

	case ScenarioBusFlood:
		// Increase latency in the RuntimeProfile to trigger a Watchdog or Trust drop
		if env.Runtime.EnvVars == nil {
			env.Runtime.EnvVars = make(map[string]string)
		}
		// Simulate a bus latency of 500ms (Critical Failure level)
		env.Runtime.EnvVars["primary_bus_latency_ms"] = "500.0"
		logging.Info("[SIM] Fault: CAN-Bus Congestion Injected")

	case ScenarioIdentityHijack:
		// Change the OS but keep the MachineID to simulate a hijacked binary
		env.Identity.OS = "unknown_hacked_os"
	}
}