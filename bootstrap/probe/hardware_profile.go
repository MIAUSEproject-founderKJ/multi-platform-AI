//bootstrap/probe/hardware_profile.go

package probe

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/math_convert"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

// ConvertHardwareFingerprint converts HardwareFingerprint to internal_environment.HardwareProfile
func ConvertFingerprintToProfile(fp HardwareFingerprint) internal_environment.HardwareProfile {
	var buses []internal_environment.BusCapability
	for bus := range fp.Buses {
		buses = append(buses, internal_environment.BusCapability{
			ID:         bus + "-bus",
			Type:       bus,
			Confidence: math_convert.FromFloat64(0.9),
			Source:     "fingerprint",
		})
	}

	return internal_environment.HardwareProfile{
		Processors: []internal_environment.Processor{
			{Type: "CPU", Count: runtime.NumCPU(), Version: 1.0},
		},
		Buses:      buses,
		HasBattery: detectBattery(),
	}
}

func detectBuses(fp *HardwareFingerprint) {
	fp.Buses = map[string]bool{}

	if len(fp.PCI) > 0 {
		fp.Buses["pci"] = true
	}

	if len(fp.MAC) > 0 {
		fp.Buses["network"] = true
	}

	// Linux-specific
	if runtime.GOOS == "linux" {
		if _, err := os.Stat("/sys/bus/i2c"); err == nil {
			fp.Buses["i2c"] = true
		}
		if _, err := os.Stat("/sys/bus/spi"); err == nil {
			fp.Buses["spi"] = true
		}
		if _, err := os.Stat("/sys/class/net/can0"); err == nil {
			fp.Buses["can"] = true
		}
	}
}

func detectBattery() bool {
	switch runtime.GOOS {
	case "linux":
		return detectBatteryLinux()
	case "darwin":
		return detectBatteryDarwin()
	case "windows":
		return detectBatteryWindows()
	default:
		return false
	}
}

func detectBatteryLinux() bool {
	entries, err := os.ReadDir("/sys/class/power_supply")
	if err != nil {
		return false
	}

	for _, e := range entries {
		if strings.HasPrefix(strings.ToLower(e.Name()), "bat") {
			return true
		}
	}
	return false
}

func detectBatteryDarwin() bool {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "Battery")
}

func detectBatteryWindows() bool {
	out, err := exec.Command("wmic", "path", "Win32_Battery", "get", "Status").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "OK")
}
