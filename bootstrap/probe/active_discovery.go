// bootstrap/probe/active_discovery.go
package probe

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

// ActiveDiscovery acts as the "Neurologist" for the machine.
// bootstrap/probe/active_discovery.go
func ActiveDiscovery(env *internal_environment.EnvConfig) (*internal_environment.EnvConfig, error) {

	logging.Info("[active_discovery] Phase 2: Active Hardware Mapping for %s", env.Platform.Final)

	switch env.Platform.Final {
	case internal_environment.PlatformComputer, internal_environment.PlatformMobile:
		populateCompute(env)
	case internal_environment.PlatformVehicle, internal_environment.PlatformRobot,
		internal_environment.PlatformIndustrial, internal_environment.PlatformEmbedded:
		populateEmbedded(env)
	default:
		logging.Warn("[PROBE] Unknown platform %s using sensor-only fallback", env.Platform.Final)
		phy, _ := discoverPhysical()
		env.Discovery.Physical = phy
		env.Discovery.Capabilities.SensorOnly = true
	}

	return env, nil
}

// populateCompute fills GPU/VRAM info for high-level devices
func populateCompute(cfg *internal_environment.EnvConfig) {
	count, totalVRAM := ProbeVRAM()
	if count > 0 {
		cfg.Hardware.Processors = append(cfg.Hardware.Processors,
			internal_environment.Processor{Type: "GPU", Count: count, Version: float64(totalVRAM)})
		if totalVRAM > 1024 {
			cfg.Discovery.Capabilities.SupportsAcceleratedCompute = true
		}
	}
}

// populateEmbedded probes layers 0–4 for embedded/vehicle/robot
func populateEmbedded(cfg *internal_environment.EnvConfig) {
	if phy, err := discoverPhysical(); err == nil {
		cfg.Discovery.Physical = phy
	} else {
		logging.Warn("physical discovery failed: %v", err)
	}
	if sig, err := discoverSignal(); err == nil {
		cfg.Discovery.Signal = sig
	}
	if nodes, err := discoverBusNodes(); err == nil {
		cfg.Discovery.Nodes = nodes
		if proto, err := discoverProtocol(nodes); err == nil {
			cfg.Discovery.Protocol = proto
			cfg.Discovery.Capabilities = resolveCapabilities(proto)
		}
	}
}

// resolveCapabilities maps protocol profile to capability descriptor
func resolveCapabilities(p internal_environment.ProtocolProfile) internal_environment.CapabilityDescriptor {
	return internal_environment.CapabilityDescriptor{
		SensorOnly:              p.ReadableRegisters > 0,
		SupportsRegisterControl: p.WritableRegisters > 0,
		SupportsGoalControl:     p.WritableRegisters > 0 && p.SupportsWatchdog && p.SupportsSafeStop,
		HasSafetyEnvelope:       p.WritableRegisters > 0 && p.SupportsWatchdog && p.SupportsSafeStop,
	}
}

// discoverPhysical probes power & voltage, cross-platform
func discoverPhysical() (internal_environment.PhysicalProfile, error) {
	phy := internal_environment.PhysicalProfile{}

	switch runtime.GOOS {
	case "linux":
		// AC power
		if data, err := os.ReadFile("/sys/class/power_supply/AC/online"); err == nil {
			phy.PowerPresent = strings.TrimSpace(string(data)) == "1"
		}
		// Battery voltage
		if data, err := os.ReadFile("/sys/class/power_supply/BAT0/voltage_now"); err == nil {
			if v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
				phy.BaseVoltage = v / 1e6
			}
		}

	case "windows":
		// On Windows, use registry or WMI defaults (simplified)
		phy.PowerPresent = true // assume AC plugged in
		phy.BaseVoltage = 12.0  // default placeholder voltage

	case "darwin":
		// macOS: assume normal AC + battery present
		phy.PowerPresent = true
		phy.BaseVoltage = 12.0

	default:
		phy.PowerPresent = true
		phy.BaseVoltage = 12.0
	}

	return phy, nil
}

// discoverSignal probes bus type and basic signal properties, cross-platform
func discoverSignal() (internal_environment.SignalProfile, error) {
	sig := internal_environment.SignalProfile{}

	switch runtime.GOOS {
	case "linux":
		out, err := exec.Command("ip", "-details", "link", "show").Output()
		if err != nil {
			return sig, err
		}
		text := string(out)
		if strings.Contains(text, "can state") {
			sig.BusType = "CAN"
			sig.StableClock = true
			for _, l := range strings.Split(text, "\n") {
				if strings.Contains(l, "bitrate") {
					parts := strings.Fields(l)
					for i, p := range parts {
						if p == "bitrate" && i+1 < len(parts) {
							if br, err := strconv.Atoi(parts[i+1]); err == nil {
								sig.BaudRate = br
							}
						}
					}
				}
			}
		}

	case "windows":
		// No CAN bus typically; fallback defaults
		sig.BusType = "ETHERNET"
		sig.StableClock = true
		sig.BaudRate = 1000_000_000 // assume 1 Gbps Ethernet

	case "darwin":
		// macOS fallback
		sig.BusType = "ETHERNET"
		sig.StableClock = true
		sig.BaudRate = 1000_000_000

	default:
		sig.BusType = "UNKNOWN"
		sig.StableClock = true
	}

	return sig, nil
}

// discoverBusNodes enumerates nodes on bus interfaces
func discoverBusNodes() ([]internal_environment.NodeDescriptor, error) {
	var nodes []internal_environment.NodeDescriptor
	if runtime.GOOS != "linux" {
		return nodes, nil
	}

	out, err := exec.Command("ls", "/sys/class/net").Output()
	if err != nil {
		return nodes, err
	}

	nodeID := 1
	for _, iface := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(iface, "can") {
			nodes = append(nodes, internal_environment.NodeDescriptor{
				NodeID:    nodeID,
				VendorID:  "UNKNOWN",
				Class:     "BusInterface",
				Heartbeat: 0,
			})
			nodeID++
		}
	}
	return nodes, nil
}

// discoverProtocol performs conservative protocol inference
func discoverProtocol(nodes []internal_environment.NodeDescriptor) (internal_environment.ProtocolProfile, error) {
	if len(nodes) == 0 {
		return internal_environment.ProtocolProfile{}, fmt.Errorf("no bus interfaces detected")
	}
	return internal_environment.ProtocolProfile{
		FirmwareVersion:   "unknown",
		WritableRegisters: 0,
		ReadableRegisters: 1,
		SupportsWatchdog:  false,
		SupportsSafeStop:  false,
	}, nil
}

// probeVRAM returns GPU count and total VRAM (MB)

func ProbeVRAM() (int, int) {
	switch runtime.GOOS {
	case "linux":
		return probeVRAMLinux()
	case "windows":
		return probeVRAMWindows()
	case "darwin":
		return probeVRAMMac()
	default:
		return 0, 0
	}
}

func probeVRAMLinux() (int, int) {
	// Nvidia GPU
	if out, err := exec.Command("nvidia-smi", "--query-gpu=memory.total", "--format=csv,noheader,nounits").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		total := 0
		for _, l := range lines {
			if v, err := strconv.Atoi(strings.TrimSpace(l)); err == nil {
				total += v
			}
		}
		return len(lines), total
	}
	// fallback: lspci
	if out, err := exec.Command("lspci").Output(); err == nil && strings.Contains(strings.ToLower(string(out)), "vga") {
		return 1, 0
	}
	return 0, 0
}

func probeVRAMWindows() (int, int) { return 0, 0 }
func probeVRAMMac() (int, int)     { return 0, 0 }
