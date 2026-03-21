//boot/probe/hardware_profile.go

package probe

import (
	"os"
	"runtime"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func collectHardwareProfile() schema.HardwareProfile {

	fp := collectHardwareFingerprint()

	var buses []schema.BusCapability

	if len(fp.PCI) > 0 {
		buses = append(buses, schema.BusCapability{
			ID:         "pci-root",
			Type:       "pci",
			Confidence: mathutil.ToFloat64(60000),
			Source:     "lspci",
		})
	}

	if len(fp.MAC) > 0 {
		buses = append(buses, schema.BusCapability{
			ID:         "ethernet",
			Type:       "network",
			Confidence: mathutil.ToFloat64(60000),
			Source:     "net-iface",
		})
	}

	return schema.HardwareProfile{
		Processors: []schema.Processor{
			{
				Type:    "CPU",
				Count:   1,
				Version: 1.0,
			},
		},
		Buses:      buses,
		HasBattery: detectBattery(),
	}
}

func detectBattery() bool {

	if runtime.GOOS != "linux" {
		return false
	}

	entries, err := os.ReadDir("/sys/class/power_supply/")
	if err != nil {
		return false
	}

	for _, e := range entries {

		name := strings.ToLower(e.Name())

		// Typical battery identifiers
		if strings.Contains(name, "bat") {
			return true
		}

		// More explicit check
		data, err := os.ReadFile("/sys/class/power_supply/" + e.Name() + "/type")
		if err == nil {
			t := strings.TrimSpace(strings.ToLower(string(data)))
			if t == "battery" {
				return true
			}
		}
	}

	return false
}
