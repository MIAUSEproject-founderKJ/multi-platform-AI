// internal/schema/env.go
// This is the "Source of Truth" that everyone can safely import.
package schema

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
)

type BootSequence struct {
	Env          *EnvConfig
	Mode         BootMode
	Attested     bool
	Capabilities CapabilitySet
	Service      ServiceType
	Entity       EntityType
	Tier         TierType
	UserSession  *UserSession
}

type EntityType uint8

const (
	EntityPersonal EntityType = iota
	EntityOrganization
	EntityStranger
	EntityTester
)

type DiscoveryDiagnostics struct {
	DiscoveryDuration time.Duration
	ProbeErrors       []string
}

func (m *MachineIdentity) BindHardware(env *EnvConfig) {
	m.Hardware = env.Hardware
}

type MachineIdentity struct {
	MachineID    string        `json:"machine_id"`
	PlatformType PlatformClass `json:"platform_type"`
	Hostname     string        `json:"hostname"`
	OS           string        `json:"os"`
	Arch         string        `json:"arch"`
	Hardware     HardwareProfile
	EntityType   EntityType
	TierType     TierType
}

type HardwareProfile struct {
	Processors []Processor     `json:"processors"`
	Buses      []BusCapability `json:"buses"`
	HasBattery bool `json:"has_battery"`
}

type BusCapability struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"` // can, usb, i2c
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
}

type Processor struct {
	Type    string  `json:"type"` // CPU, GPU, TPU
	Count   int     `json:"count"`
	Version float64 `json:"version"`
}


type EnvConfig struct {
	SchemaVersion int                `json:"schema_version"`
	Discovery     DiscoveryProfile   `json:"discovery_profile"`
	GeneratedAt   time.Time          `json:"generated_at"`
	Identity      MachineIdentity    `json:"identity"`
	Hardware      HardwareProfile    `json:"hardware"`
	Platform      PlatformResolution `json:"platform"`
	Attestation   EnvAttestation     `json:"attestation"`
	EntityType    EntityType         `json:"entity_type"`
	TierType      TierType           `json:"tier_type"`
}

// EnvAttestation defines the cryptographic seal of the environment
type EnvAttestation struct {
	Locked        bool          `json:"locked"`
	PlatformClass PlatformClass `json:"platform_class,omitempty"`
	Valid         bool          `json:"valid"`
	Level         BootTrust     `json:"level"` // "strong" | "weak" | "invalid"
	EnvHash       string        `json:"env_hash"`
	SessionToken  string        `json:"session_token,omitempty"`
}

type BootTrust uint8

const (
	TrustInvalid BootTrust = iota
	TrustWeak
	TrustStrong
)

type IdentityProfile struct {
	MachineID    string
	MachineName  string
	OS           string
	Architecture string
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
