// bootstrap/resolver/capabilities_resolver.go
package bootstrap_resolver

import (
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

func DeviceCapabilitiesResolver(env *internal_environment.EnvConfig, capSet internal_verification.CapabilitySet) (*internal_environment.CapabilityDescriptor, error) {
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
	if capSet&internal_verification.CapCANBus != 0 {
		caps.SensorOnly = false
	}

	if capSet&internal_verification.CapSafetyCritical != 0 {
		caps.HasSafetyEnvelope = true
	}

	if capSet&internal_verification.CapLocalStorage != 0 {
		caps.SupportsRegisterControl = true
	}

	// --------------------------------
	// 3. Platform-specific logic
	// --------------------------------
	switch env.Platform.Final {

	case internal_environment.PlatformRobot, internal_environment.PlatformVehicle:
		caps.SupportsGoalControl =
			caps.SupportsRegisterControl &&
				caps.HasSafetyEnvelope

	case internal_environment.PlatformEmbedded:
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

func extractPlatformConfidence(env *internal_environment.EnvConfig) float64 {
	for _, c := range env.Platform.Candidates {
		if c.Type == env.Platform.Final {
			return c.Confidence.Float64()
		}
	}
	return 0.5 // fallback
}
