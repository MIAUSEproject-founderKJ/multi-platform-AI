//boot/probe/capability_mapper.go

package probe

import (
	"runtime"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func BuildCapabilitySet(env *schema.EnvConfig, fp HardwareFingerprint) schema.CapabilitySet {

	var caps schema.CapabilitySet

	hw := env.Hardware

	// --------------------------------------------------
	// CORE SYSTEM CAPABILITIES (baseline)
	// --------------------------------------------------

	if runtime.GOOS != "" {
		caps.Add(schema.CapNetwork)
		caps.Add(schema.CapLocalStorage)
	}

	if len(hw.Processors) > 0 {
		caps.Add(schema.CapLocalStorage)
	}

	// --------------------------------------------------
	// SECURITY CAPABILITIES
	// --------------------------------------------------

	if fp.TPM != "" {
		caps.Add(schema.CapSecureEnclave)
	}

	// Heuristic fallback (modern OS = soft enclave)
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		caps.Add(schema.CapSecureEnclave)
	}

	// --------------------------------------------------
	// HARDWARE / PLATFORM CAPABILITIES
	// --------------------------------------------------

	for _, b := range hw.Buses {

		switch b.Type {

		case "network":
			caps.Add(schema.CapNetwork)

		case "pci":
			caps.Add(schema.CapLocalStorage)

		case "usb":
			caps.Add(schema.CapLocalStorage)

		case "can":
			caps.Add(schema.CapCANBus)
			caps.Add(schema.CapSafetyCritical)

		case "industrial":
			caps.Add(schema.CapIndustrialIO)
			caps.Add(schema.CapSafetyCritical)
		}
	}

	// --------------------------------------------------
	// DEVICE CLASS INFERENCE (HIGH IMPACT)
	// --------------------------------------------------

	// Laptop / mobile indicator
	if hw.HasBattery {
		caps.Add(schema.CapBiometric) // assume user device
	}

	// GPU / compute acceleration (if you add detection)
	if hasGPU(hw) {
		caps.Add(schema.CapPersistentCloudLink)
	}

	// --------------------------------------------------
	// NEGATIVE SIGNALS (CRITICAL FOR STABILITY)
	// --------------------------------------------------

	if !hasBus(hw, "can") && !hasBus(hw, "industrial") {
		// strongly implies non-safety-critical environment
		caps.Remove(schema.CapSafetyCritical)
	}

	return caps
}

func hasBus(fp HardwareFingerprint, bus string) bool {
	return fp.Buses[bus]
}

func hasGPU(hw schema.HardwareProfile) bool {
	// extend later with real detection
	return false
}
