// core/bootstrap.go

package core

import (
    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ResolveRuntimeContext bridges the Hardware Probe (schema.EnvConfig) 
// to the Module Policy (core.RuntimeContext).
func ResolveRuntimeContext(env *schema.EnvConfig, bootMode string) RuntimeContext {
    
    // 1. Map Physical Hardware to Abstract Capabilities
    // This decouples modules from specific drivers.
    caps := make(map[Capability]bool)

    // Hardware-derived capabilities
    if env.Discovery.Signal.BusType == "CAN" {
        caps[CapCANBus] = true
    }
    if env.Discovery.Physical.PowerPresent {
        caps[CapActuators] = true // Assume actuators if high-voltage rail exists
    }
    
    // Schema-derived capabilities
    if env.Discovery.Capabilities.HasSafetyEnvelope {
        caps[CapSafetyCritical] = true
    }
    if env.Discovery.Capabilities.SupportsGoalControl {
        caps[CapHighFreqSensor] = true
    }
    
    // Platform-class derivatives (Legacy compatibility)
    if env.Platform.Final == "Workstation" || env.Platform.Final == "Laptop" {
        caps[CapFileSystem] = true
        caps[CapMicrophone] = true
        caps[CapGPU] = true // Simplified assumption, ideally probed
    }

    // 2. Resolve Entity & Policy (Who owns this?)
    // This comes from the cryptographic passport, not the hardware.
    entityType := "Stranger"
    permissions := make(map[string]bool)
    
    if env.Attestation.Valid {
        entityType = "VerifiedNode"
        permissions["TRUSTED_BOOT"] = true
        
        // If we are fully verified, allow autonomous execution
        if env.Platform.Mode == "AUTONOMOUS" {
            permissions["AUTONOMOUS_EXECUTION"] = true
        }
    }

    // 3. Construct the Immutable Context
    return RuntimeContext{
        Platform: PlatformProfile{
            Name:         string(env.Platform.Final),
            Capabilities: caps,
        },
        Identity: IdentityProfile{
            Entity: entityType,
        },
        Policy: PolicyProfile{
            Permissions: permissions,
        },
        Boot: BootProfile{
            Type: bootMode,
        },
        // Service profile would typically be injected by the HMI or Cloud config
        Service: ServiceProfile{ Name: "DefaultAgent" }, 
    }
}