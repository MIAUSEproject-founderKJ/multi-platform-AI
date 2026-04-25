// boot/detect_cap.go
package boot

import (
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

func DetectDeviceCapabilities(env *schema_system.EnvConfig, capSet schema_security.CapabilitySet) (*schema_system.CapabilityDescriptor, error) {
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
	if capSet&schema_security.CapCANBus != 0 {
		caps.SensorOnly = false
	}

	if capSet&schema_security.CapSafetyCritical != 0 {
		caps.HasSafetyEnvelope = true
	}

	if capSet&schema_security.CapLocalStorage != 0 {
		caps.SupportsRegisterControl = true
	}

	// --------------------------------
	// 3. Platform-specific logic
	// --------------------------------
	switch env.Platform.Final {

	case schema_system.PlatformRobot, schema_system.PlatformVehicle:
		caps.SupportsGoalControl =
			caps.SupportsRegisterControl &&
				caps.HasSafetyEnvelope

	case schema_system.PlatformEmbedded:
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

func extractPlatformConfidence(env *schema_system.EnvConfig) float64 {
	for _, c := range env.Platform.Candidates {
		if c.Type == env.Platform.Final {
			return c.Confidence.Float64()
		}
	}
	return 0.5 // fallback
}
