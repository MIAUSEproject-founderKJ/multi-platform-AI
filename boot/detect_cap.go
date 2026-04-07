// boot/detect_cap.go
package boot

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func DetectDeviceCapabilities(env *schema.EnvConfig, capSet schema.CapabilitySet) (*schema.CapabilityDescriptor, error) {
	caps := env.Discovery.Capabilities // start from discovery

	// --------------------------------
	// 1. Hardware enrichment
	// --------------------------------
	for _, p := range env.Hardware.Processors {
		if p.Type == "GPU" && p.Count > 0 {
			caps.SupportsAcceleratedCompute = true
		}
	}

	// --------------------------------
	// 2. Platform / Boot capability merge
	// --------------------------------
	if capSet&schema.CapCANBus != 0 {
		caps.SensorOnly = false
	}

	if capSet&schema.CapSafetyCritical != 0 {
		caps.HasSafetyEnvelope = true
	}

	if capSet&schema.CapLocalStorage != 0 {
		caps.SupportsRegisterControl = true
	}

	// --------------------------------
	// 3. Platform-specific logic
	// --------------------------------
	switch env.Platform.Final {

	case schema.PlatformRobot, schema.PlatformVehicle:
		caps.SupportsGoalControl =
			caps.SupportsRegisterControl &&
				caps.HasSafetyEnvelope

	case schema.PlatformEmbedded:
		caps.SensorOnly = !caps.SupportsRegisterControl
	}

	// --------------------------------
	// 4. Safety fallback
	// --------------------------------
	if !caps.SupportsGoalControl && !caps.SupportsRegisterControl {
		caps.SensorOnly = true
	}

	return &caps, nil
}

func extractPlatformConfidence(env *schema.EnvConfig) float64 {
	for _, c := range env.Platform.Candidates {
		if c.Type == env.Platform.Final {
			return c.Confidence.Float64()
		}
	}
	return 0.5 // fallback
}
