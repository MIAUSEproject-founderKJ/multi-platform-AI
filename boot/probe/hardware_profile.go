//boot/probe/hardware_profile.go

package probe

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func collectHardwareProfile() schema.HardwareProfile {

	fp := collectHardwareFingerprint()

	var buses []schema.BusCapability

	if len(fp.PCI) > 0 {
		buses = append(buses, schema.BusCapability{
			ID: "pci-root",
			Type: "pci",
			Confidence: 60000,
			Source: "lspci",
		})
	}

	if len(fp.MAC) > 0 {
		buses = append(buses, schema.BusCapability{
			ID: "ethernet",
			Type: "network",
			Confidence: 60000,
			Source: "net-iface",
		})
	}

	return schema.HardwareProfile{
		Processors: []schema.Processor{
			{
				Type: "CPU",
				Count: 1,
				Version: 1.0,
			},
		},
		Buses: buses,
		HasBattery: detectBattery(),
	}
}