// core/platform/probe/active.go
package probe

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ActiveDiscovery acts as the "Neurologist" for the machine.
func ActiveDiscovery(id *HardwareIdentity) (*schema.EnvConfig, error) {
	logging.Info("[PROBE] Phase 2: Starting Active Hardware Mapping for %s...", id.PlatformType)

	// Initialize the config with the Passport data
	config := &schema.EnvConfig{
		SchemaVersion: schema.CurrentVersion,
		GeneratedAt:   time.Now(),
		Identity: schema.MachineIdentity{
			MachineID:   id.InstanceID,
			OS:          id.OS,
			Arch:        id.Arch,
		},

		Platform: schema.PlatformResolution{
		Final:      id.PlatformType,
		Source:     "active_probe_v1",
		ResolvedAt: time.Now(),
		Locked:     false,
	},
}

	// ------------------------------------------
	// PATH A: High-Level Compute 
	// ------------------------------------------
switch id.PlatformType {
case schema.PlatformComputer, schema.PlatformLaptop, schema.PlatformTablet, schema.PlatformMobile:
    // High-level compute path
    gpuCount, totalVRAM := probeVRAM()
    if gpuCount > 0 {
        config.Hardware.Processors = append(config.Hardware.Processors,
            schema.Processor{Type: "GPU", Count: gpuCount, Version: float64(totalVRAM)})
    }
    if gpuCount > 0 && totalVRAM > 1024 {
        config.Discovery.Capabilities.SupportsAcceleratedCompute = true
    }

case schema.PlatformVehicle, schema.PlatformDrone, schema.PlatformRobot, schema.PlatformIndustrial, schema.PlatformEmbedded, schema.PlatformGamePad:
    // Deep-layered probe
    phy, err := discoverPhysical()
    if err != nil {
        logging.Warn("[PROBE] Physical layer failure: %v", err)
    } else {
        config.Discovery.Physical = phy
    }

    sig, err := discoverSignal()
    if err != nil {
        logging.Warn("[PROBE] Signal layer failure: %v", err)
    } else {
        config.Discovery.Signal = sig
    }

    nodes, err := discoverBusNodes()
    if err != nil {
        logging.Warn("[PROBE] Bus enumeration failed: %v", err)
    } else {
        config.Discovery.Nodes = nodes
    }

    proto, err := discoverProtocol(nodes)
    if err != nil {
        logging.Warn("[PROBE] Protocol discovery failed: %v", err)
    } else {
        config.Discovery.Protocol = proto
    }

    config.Discovery.Capabilities = resolveCapabilities(config.Discovery.Protocol)

default:
    logging.Warn("[PROBE] Unknown platform type %s, using sensor-only default", id.PlatformType)
    phy, _ := discoverPhysical()
    config.Discovery.Physical = phy
    config.Discovery.Capabilities.SensorOnly = true
}

// Final return
return config, nil
}


func resolveCapabilities(p schema.ProtocolProfile) schema.CapabilityDescriptor {
	c := schema.CapabilityDescriptor{}

	if p.ReadableRegisters > 0 {
		c.SensorOnly = true
	}

	if p.WritableRegisters > 0 {
		c.SupportsRegisterControl = true
	}

	if p.WritableRegisters > 0 &&
		p.SupportsWatchdog &&
		p.SupportsSafeStop {
		c.SupportsGoalControl = true
		c.HasSafetyEnvelope = true
	}

	return c
}

func discoverPhysical() (schema.PhysicalProfile, error) {

    phy := schema.PhysicalProfile{}

    if runtime.GOOS == "linux" {

        if data, err := os.ReadFile("/sys/class/power_supply/AC/online"); err == nil {
            phy.PowerPresent = strings.TrimSpace(string(data)) == "1"
        }

        if data, err := os.ReadFile("/sys/class/power_supply/BAT0/voltage_now"); err == nil {
            if v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
                phy.BaseVoltage = v / 1e6
            }
        }
    }

    return phy, nil
}


func discoverSignal() (schema.SignalProfile, error) {

	sig := schema.SignalProfile{}

	if runtime.GOOS == "linux" {

		
		out, err := exec.Command("ip", "-details", "link", "show").Output()
		if err != nil {
			return sig, err
		}

		text := string(out)

		if strings.Contains(text, "can state") {
			sig.BusType = "CAN"
			sig.StableClock = true

			// Attempt to parse bitrate
			lines := strings.Split(text, "\n")
			for _, l := range lines {
				if strings.Contains(l, "bitrate") {
					parts := strings.Fields(l)
					for i, p := range parts {
						if p == "bitrate" && i+1 < len(parts) {
							br, _ := strconv.Atoi(parts[i+1])
							sig.BaudRate = br
						}
					}
				}
			}
		}
	}

	return sig, nil
}


func discoverBusNodes() ([]schema.NodeDescriptor, error) {

	var nodes []schema.NodeDescriptor

	if runtime.GOOS == "linux" {

		out, err := exec.Command("ls", "/sys/class/net").Output()
		if err != nil {
			return nil, err
		}

		interfaces := strings.Split(string(out), "\n")

		nodeID := 1
		for _, iface := range interfaces {
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
	}

	return nodes, nil
}


func discoverProtocol(nodes []schema.NodeDescriptor) (schema.ProtocolProfile, error) {

	if len(nodes) == 0 {
		return schema.ProtocolProfile{}, fmt.Errorf("no bus interfaces detected")
	}

	// Conservative inference
	return schema.ProtocolProfile{
		FirmwareVersion:   "unknown",
		WritableRegisters: 0,
		ReadableRegisters: 1,
		SupportsWatchdog:  false,
		SupportsSafeStop:  false,
	}, nil
}



func probeVRAM() (int, int) {
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

	// Try NVIDIA first
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=memory.total",
		"--format=csv,noheader,nounits")

	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		total := 0
		for _, l := range lines {
			v, err := strconv.Atoi(strings.TrimSpace(l))
			if err == nil {
				total += v
			}
		}
		return len(lines), total
	}

	// Fallback: lspci detect GPU presence
	out, err = exec.Command("lspci").Output()
	if err == nil && strings.Contains(strings.ToLower(string(out)), "vga") {
		return 1, 0
	}

	return 0, 0
}
