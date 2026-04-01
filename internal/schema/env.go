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

// Use EntityType (uint8) internally for speed and clarity. Use TierType (string) externally for readability and compatibility.
type EntityType uint8

const (
	EntityPersonal EntityType = iota
	EntityOrganization
	EntityStranger
	EntityTester
)

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
	TierType     TierType //Use TierType (string) externally for readability and compatibility. Use EntityType (uint8) internally for speed and clarity.
	Password     string   `json:"password,omitempty"` // For authentication during cold boot. Should be securely handled and not stored in plaintext in production.
}

type HardwareProfile struct {
	Processors []Processor     `json:"processors"`
	Buses      []BusCapability `json:"buses"`
	HasBattery bool            `json:"has_battery"`
}

type BusCapability struct {
	ID         string       `json:"id"`
	Type       string       `json:"type"` // can, usb, i2c
	Confidence mathutil.Q16 `json:"confidence"`
	Source     string       `json:"source"`
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
	Physical          PhysicalProfile      `json:"physical"`
	Signal            SignalProfile        `json:"signal"`
	Nodes             []NodeDescriptor     `json:"nodes"`
	Protocol          ProtocolProfile      `json:"protocol"`
	Capabilities      CapabilityDescriptor `json:"capabilities"`
	DiscoveryDuration time.Duration        `json:"discovery_duration"`
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
