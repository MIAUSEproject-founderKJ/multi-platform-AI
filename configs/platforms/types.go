//MIAUSEproject-founderKJ/multi-platform-AI/configs/platforms/types.go

package platforms

import "time"

// PlatformClass defines the categorical identity of the device.
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

type MachineIdentity struct {
    MachineName string `json:"machine_name"`
    OS          string `json:"os"`
    Arch        string `json:"arch"`
}

type HardwareProfile struct {
    Processors []Processor     `json:"processors"`
    Buses      []BusCapability `json:"buses"`
    Battery    BatteryStatus   `json:"battery"`
}

type BusCapability struct {
    ID         string      `json:"id"`   
    Type       string      `json:"type"` // can, usb, i2c
    Confidence uint16      `json:"confidence"` // Q16 format
    Source     string      `json:"source"`
}

type Processor struct {
    Type    string  `json:"type"`    // CPU, GPU, TPU
    Count   int     `json:"count"`
    Version float64 `json:"version"`
}