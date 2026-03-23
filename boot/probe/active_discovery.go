// boot/probe/active_discovery.go
package probe

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ActiveDiscovery acts as the "Neurologist" for the machine.
func ActiveDiscovery(env *schema.EnvConfig) (*schema.EnvConfig, error) {

	logging.Info(
		"[PROBE] Phase 2: Active Hardware Mapping for %s",
		env.Platform.Final,
	)

	env.GeneratedAt = time.Now()

	switch env.Platform.Final {

	case schema.PlatformComputer,
		schema.PlatformLaptop,
		schema.PlatformTablet,
		schema.PlatformMobile:

		populateCompute(env)

	case schema.PlatformVehicle,
		schema.PlatformDrone,
		schema.PlatformRobot,
		schema.PlatformIndustrial,
		schema.PlatformEmbedded,
		schema.PlatformGamePad:

		populateEmbedded(env)

	default:

		logging.Warn(
			"[PROBE] Unknown platform %s using sensor-only fallback",
			env.Platform.Final,
		)

		phy, _ := discoverPhysical()
		env.Discovery.Physical = phy
		env.Discovery.Capabilities.SensorOnly = true
	}

	return env, nil
}

// populateCompute fills GPU/VRAM info for high-level devices
func populateCompute(cfg *schema.EnvConfig) {
	count, totalVRAM := ProbeVRAM()
	if count > 0 {
		cfg.Hardware.Processors = append(cfg.Hardware.Processors,
			schema.Processor{Type: "GPU", Count: count, Version: float64(totalVRAM)})
		if totalVRAM > 1024 {
			cfg.Discovery.Capabilities.SupportsAcceleratedCompute = true
		}
	}
}

// populateEmbedded probes layers 0–4 for embedded/vehicle/robot
func populateEmbedded(cfg *schema.EnvConfig) {
	if phy, err := discoverPhysical(); err == nil {
		cfg.Discovery.Physical = phy
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
func resolveCapabilities(p schema.ProtocolProfile) schema.CapabilityDescriptor {
	return schema.CapabilityDescriptor{
		SensorOnly:              p.ReadableRegisters > 0,
		SupportsRegisterControl: p.WritableRegisters > 0,
		SupportsGoalControl:     p.WritableRegisters > 0 && p.SupportsWatchdog && p.SupportsSafeStop,
		HasSafetyEnvelope:       p.WritableRegisters > 0 && p.SupportsWatchdog && p.SupportsSafeStop,
	}
}

// discoverPhysical probes power & voltage
func discoverPhysical() (schema.PhysicalProfile, error) {
	phy := schema.PhysicalProfile{}

	if runtime.GOOS != "linux" {
		return phy, nil
	}

	if data, err := os.ReadFile("/sys/class/power_supply/AC/online"); err == nil {
		phy.PowerPresent = strings.TrimSpace(string(data)) == "1"
	}
	if data, err := os.ReadFile("/sys/class/power_supply/BAT0/voltage_now"); err == nil {
		if v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
			phy.BaseVoltage = v / 1e6
		}
	}

	return phy, nil
}

// discoverSignal probes bus type and basic signal properties
func discoverSignal() (schema.SignalProfile, error) {
	sig := schema.SignalProfile{}
	if runtime.GOOS != "linux" {
		return sig, nil
	}

	out, err := exec.Command("ip", "-details", "link", "show").Output()
	if err != nil {
		return sig, err
	}
	text := string(out)
	if !strings.Contains(text, "can state") {
		return sig, nil
	}

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

	return sig, nil
}

// discoverBusNodes enumerates nodes on bus interfaces
func discoverBusNodes() ([]schema.NodeDescriptor, error) {
	var nodes []schema.NodeDescriptor
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
			nodes = append(nodes, schema.NodeDescriptor{
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
func discoverProtocol(nodes []schema.NodeDescriptor) (schema.ProtocolProfile, error) {
	if len(nodes) == 0 {
		return schema.ProtocolProfile{}, fmt.Errorf("no bus interfaces detected")
	}
	return schema.ProtocolProfile{
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
