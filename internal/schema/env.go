// internal/schema/env.go
// This is the "Source of Truth" that everyone can safely import.
package schema

import (
	"time"
)

// PlatformClass defines the type of hardware (Vehicle, Drone, etc.)
type PlatformClass string

const (
	PlatformComputer   PlatformClass = "computer"
	PlatformLaptop     PlatformClass = "laptop"
	PlatformMobile     PlatformClass = "mobile"
	PlatformTablet     PlatformClass = "tablet"
	PlatformRobot      PlatformClass = "robotic"
	PlatformVehicle    PlatformClass = "vehicle"
	PlatformDrone      PlatformClass = "drone"
	PlatformIndustrial PlatformClass = "industrial"
	PlatformEmbedded   PlatformClass = "embedded"
	PlatformGamePad    PlatformClass = "gamepad"
)

type PlatformProfile struct {
	Final      PlatformClass   `json:"final_class"`
	Candidates []PlatformScore `json:"candidates"`
	Source     string          `json:"source"`
	ResolvedAt time.Time       `json:"resolved_at"`
	Locked     bool            `json:"locked"`
}

// PlatformScore tracks the heuristic weight for a specific platform type.
type PlatformScore struct {
	Type       PlatformClass `json:"type"`
	Score      float64       `json:"raw_score"`
	MaxScore   float64       `json:"max_score"`  // Potential maximum for normalization
	Confidence uint16        `json:"confidence"` // Normalized Q16 (0-65535)
	Signals    []string      `json:"signals"`    // Evidence found (e.g., "CAN_BUS_PRESENT")
}

type BootSequence struct {
	PlatformID string
	TrustScore float64
	IsVerified bool
	Mode       string
	UserRole   string
	EnvConfig  *schema.EnvConfig
}

type EnvConfig struct {
	SchemaVersion int                `json:"schema_version"`
	Discovery     DiscoveryProfile	 `json:"discovery_profile"`
	GeneratedAt   time.Time          `json:"generated_at"`
	Identity      MachineIdentity    `json:"identity"`
	Hardware      HardwareProfile    `json:"hardware"`
	Platform      PlatformResolution `json:"platform"`
	Attestation   EnvAttestation     `json:"attestation"`
}


type MachineIdentity struct {
	MachineID   string `json:"machine_id"`
	MachineName string `json:"machine_name"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
}

type HardwareProfile struct {
	Processors []Processor     `json:"processors"`
	Buses      []BusCapability `json:"buses"`
	// Simplified Battery for the check
	HasBattery bool `json:"has_battery"`
}

type BusCapability struct {
	ID         string `json:"id"`
	Type       string `json:"type"` // can, usb, i2c
	Confidence uint16 `json:"confidence"`
	Source     string `json:"source"`
}

type Processor struct {
	Type    string  `json:"type"` // CPU, GPU, TPU
	Count   int     `json:"count"`
	Version float64 `json:"version"`
}

// PlatformResolution is the finalized identity of the environment.
type PlatformResolution struct {
	Candidates []PlatformScore `json:"candidates"`
	Final      PlatformClass   `json:"final"`
	Locked     bool            `json:"locked"`
	Source     string          `json:"source"` // e.g., "heuristic_v1" or "manual_override"
	ResolvedAt time.Time       `json:"resolved_at"`
}



// EnvAttestation defines the cryptographic seal of the environment
type EnvAttestation struct {
	Valid        bool   `json:"valid"`
	Level        string `json:"level"` // "strong" | "weak" | "invalid"
	EnvHash      string `json:"env_hash"`
	SessionToken string `json:"session_token,omitempty"`
}

type IdentityProfile struct {
	MachineID   string
	MachineName string
	OS          string
}

type BusEntry struct {
	ID         string
	Type       string
	Confidence uint16 // Q16 format
}



type DiscoveryProfile struct {
	Physical     PhysicalProfile      `json:"physical"`
	Signal       SignalProfile        `json:"signal"`
	Nodes        []NodeDescriptor     `json:"nodes"`
	Protocol     ProtocolProfile      `json:"protocol"`
	Capabilities CapabilityDescriptor `json:"capabilities"`
}

type PhysicalProfile struct {
	PowerPresent bool    `json:"power_present"`
	BaseVoltage  float64 `json:"base_voltage"`
}

type SignalProfile struct {
	BusType     string  `json:"bus_type"`
	BaudRate    int     `json:"baud_rate"`
	NoiseLevel  float64 `json:"noise_level"`
	StableClock bool    `json:"stable_clock"`
}

type NodeDescriptor struct {
	NodeID    int    `json:"node_id"`
	VendorID  string `json:"vendor_id"`
	Class     string `json:"class"` // e.g., "Actuator", "Sensor"
	Heartbeat int    `json:"heartbeat_ms"`
}

type ProtocolProfile struct {
	FirmwareVersion   string `json:"firmware_version"`
	WritableRegisters int    `json:"writable_registers"`
	ReadableRegisters int    `json:"readable_registers"`
	SupportsWatchdog  bool   `json:"supports_watchdog"`
	SupportsSafeStop  bool   `json:"supports_safe_stop"`
}

type CapabilityDescriptor struct {
	SupportsGoalControl     bool `json:"supports_goal_control"`     // AI says "Go 10m"
	SupportsRegisterControl bool `json:"supports_register_control"` // AI says "Motor Voltage 5V"
	SensorOnly              bool `json:"sensor_only"`
	HasSafetyEnvelope       bool `json:"has_safety_envelope"`
}