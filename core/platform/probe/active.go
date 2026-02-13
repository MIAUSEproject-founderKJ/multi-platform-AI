// core/platform/probe/active.go
package probe

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ActiveDiscovery acts as the "Neurologist" for the machine.
func ActiveDiscovery(id *HardwareIdentity) (*schema.EnvConfig, error) {
	logging.Info("[PROBE] Phase 2: Starting Active Hardware Mapping for %s...", id.PlatformType)

	// Initialize the config with the Passport data
	config := &schema.EnvConfig{
		PlatformID:    id.InstanceID,
		SchemaVersion: schema.CurrentVersion,
		Platform:      schema.PlatformProfile{Final: id.PlatformType}, 
		Capabilities:  make(map[string]bool),
	}

	// ------------------------------------------
	// PATH A: High-Level Compute (Laptop/Server)
	// ------------------------------------------
	if id.PlatformType == "Workstation" || id.PlatformType == "Laptop" {
		config.VRAMTotal = probeVRAM()
		config.Capabilities["GPU_ACCELERATION"] = config.VRAMTotal > 1024
		return config, nil
	}

	// ------------------------------------------
	// PATH B: Embedded Control (Vehicle/Robot)
	// ------------------------------------------
	// Only run the deep layered probe on hardware that supports it
	
	// Layer 0: Physical
	phy, err := discoverPhysical()
	if err != nil {
		return nil, fmt.Errorf("layer 0 failure: %w", err)
	}
	config.Discovery.Physical = phy

	// Layer 1: Signal
	sig, err := discoverSignal()
	if err != nil {
		return nil, fmt.Errorf("layer 1 failure: %w", err)
	}
	config.Discovery.Signal = sig

	// Layer 2: Bus Enumeration
	nodes, err := discoverBusNodes()
	if err != nil {
		logging.Warn("[PROBE] Layer 2 Warning: %v", err)
	}
	config.Discovery.Nodes = nodes

	// Layer 3: Protocol
	proto, err := discoverProtocol(nodes)
	if err != nil {
		logging.Warn("[PROBE] Layer 3 Warning: %v", err)
	}
	config.Discovery.Protocol = proto

	// Layer 4: Capability Resolution
	caps := resolveCapabilities(proto)
	config.Discovery.Capabilities = caps

	// Map findings to high-level Capabilities for the Kernel
	if caps.SupportsGoalControl {
		config.Capabilities["AUTONOMOUS_ACTUATION"] = true
	}

	logging.Info("[PROBE] Layered Discovery Complete. Safety Envelope: %v", caps.HasSafetyEnvelope)
	return config, nil
}
func resolveCapabilities(p schema.ProtocolProfile) schema.CapabilityDescriptor {
	c := schema.CapabilityDescriptor{}

	// Register-level control
	if p.WritableRegisters > 0 {
		c.SupportsRegisterControl = true
	}

	// Supervisory control requires safety primitives
	if p.SupportsWatchdog && p.SupportsSafeStop && p.WritableRegisters > 0 {
		c.SupportsGoalControl = true
		c.HasSafetyEnvelope = true
	}

	// Sensor-only classification
	if p.WritableRegisters == 0 && p.ReadableRegisters > 0 {
		c.SensorOnly = true
	}

	return c
}

func discoverPhysical() (schema.PhysicalProfile, error) {
	// hardware-specific checks
	return schema.PhysicalProfile{
		PowerPresent: true,
		BaseVoltage:  12.4,
	}, nil
}

func discoverSignal() (schema.SignalProfile, error) {
	return schema.SignalProfile{
		BusType:     "CAN",
		BaudRate:    500000,
		NoiseLevel:  0.05,
		StableClock: true,
	}, nil
}

func discoverBusNodes() ([]schema.NodeDescriptor, error) {
	return []schema.NodeDescriptor{
		{NodeID: 1, VendorID: "ACME", Class: "Actuator", Heartbeat: 100},
	}, nil
}

func discoverProtocol(nodes []schema.NodeDescriptor) (schema.ProtocolProfile, error) {
	if len(nodes) == 0 {
		return schema.ProtocolProfile{}, fmt.Errorf("no nodes to interrogate")
	}

	return schema.ProtocolProfile{
		FirmwareVersion:   "1.2.3",
		WritableRegisters: 16,
		ReadableRegisters: 32,
		SupportsWatchdog:  true,
		SupportsSafeStop:  true,
	}, nil
}
