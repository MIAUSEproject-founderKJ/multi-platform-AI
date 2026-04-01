//internal/schema/device_platform.go

package schema

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
)

type DeviceClass string

const (
	DeviceComputer   DeviceClass = "computer"
	DeviceMobile     DeviceClass = "mobile"
	DeviceEmbedded   DeviceClass = "embedded"
	DeviceIndustrial DeviceClass = "industrial"
	DeviceVehicle    DeviceClass = "vehicle"
	DeviceRobot      DeviceClass = "robot"
)

type FormFactor string

const (
	FormDesktop  FormFactor = "desktop"
	FormLaptop   FormFactor = "laptop"
	FormTablet   FormFactor = "tablet"
	FormPhone    FormFactor = "phone"
	FormHandheld FormFactor = "handheld"
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
	Confidence mathutil.Q16
	Q16        mathutil.Q16
}

type Signal struct {
	Name       string
	Value      float64
	Confidence mathutil.Q16
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
