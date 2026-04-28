//internal/schema/system/hardware_profile.go

package internal_environment

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/math_convert"
)

type HardwareProfile struct {
	Processors []Processor     `json:"processors"`
	Buses      []BusCapability `json:"buses"`
	HasBattery bool            `json:"has_battery"`
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

type ProtocolProfile struct {
	FirmwareVersion   string `json:"firmware_version"`
	WritableRegisters int    `json:"writable_registers"`
	ReadableRegisters int    `json:"readable_registers"`
	SupportsWatchdog  bool   `json:"supports_watchdog"`
	SupportsSafeStop  bool   `json:"supports_safe_stop"`
}

type CapabilityDescriptor struct {
	SensorOnly                 bool
	SupportsRegisterControl    bool
	SupportsGoalControl        bool
	HasSafetyEnvelope          bool
	SupportsAcceleratedCompute bool
}

type NodeDescriptor struct {
	NodeID    int    `json:"node_id"`
	VendorID  string `json:"vendor_id"`
	Class     string `json:"class"` // e.g., "Actuator", "Sensor"
	Heartbeat int    `json:"heartbeat_ms"`
}

type BusEntry struct {
	ID         string
	Type       string
	Confidence uint16 // Q16 format
}

type BusCapability struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"` // can, usb, i2c
	Confidence math_convert.Q16 `json:"confidence"`
	Source     string           `json:"source"`
}

type Processor struct {
	Type    string  `json:"type"` // CPU, GPU, TPU
	Count   int     `json:"count"`
	Version float64 `json:"version"`
}

type DeviceClass string

const (
	DeviceComputer   DeviceClass = "computer"
	DeviceMobile     DeviceClass = "mobile"
	DeviceEmbedded   DeviceClass = "embedded"
	DeviceIndustrial DeviceClass = "industrial"
	DeviceVehicle    DeviceClass = "vehicle"
	DeviceRobot      DeviceClass = "robot"
)

type CapabilityTag string

const (
	TagDrone      CapabilityTag = "drone"
	TagGamepad    CapabilityTag = "gamepad"
	TagAutomotive CapabilityTag = "automotive"
)

// PlatformClass defines the type of hardware (Vehicle, Drone, etc.)
type PlatformClass string

const (
	PlatformComputer   PlatformClass = "computer"
	PlatformMobile     PlatformClass = "mobile"
	PlatformEmbedded   PlatformClass = "embedded"
	PlatformIndustrial PlatformClass = "industrial"
	PlatformVehicle    PlatformClass = "vehicle"
	PlatformRobot      PlatformClass = "robot"
	PlatformUnknown    PlatformClass = "unknown"
)

type PlatformProfile struct {
	Class        DeviceClass
	Form         FormFactor
	Capabilities []CapabilityTag
}

// PlatformScore tracks the heuristic weight for a specific platform type.

type PlatformScore struct {
	Type       PlatformClass
	Profile    PlatformProfile // NEW
	Signals    []Signal
	Score      float64
	MaxScore   float64
	Confidence math_convert.Q16
	Q16        math_convert.Q16
}

type Signal struct {
	Name       string
	Value      float64
	Confidence math_convert.Q16
	Weight     float64
	Source     string
}

// PlatformResolution is the finalized identity of the environment.
type PlatformResolution struct {
	Candidates []PlatformScore `json:"candidates"`
	Final      PlatformClass   `json:"final"`
	Locked     bool            `json:"locked"`
	Source     string          `json:"source"` // e.g., "heuristic_v1" or "manual_override"
	ResolvedAt time.Time       `json:"resolved_at"`
}

type Capability uint64
type CapabilitySet uint64

const (
	CapDisplay Capability = 1 << iota
	CapKeyboard
	CapTouch
	CapMicrophone
	CapSpeaker
	CapCamera
	CapGPU
	CapSecureEnclave
	CapNetwork
	CapCANBus
	CapBiometric
	CapHighFreqSensor
	CapFileSystem
	CapSafetyCritical
	CapPersistentCloudLink
	CapIndustrialIO
	CapLocalStorage
)
