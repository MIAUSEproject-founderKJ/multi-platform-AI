//internal/schema/env.go
//This is the "Source of Truth" that everyone can safely import.
package schema

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

// BootSequence represents the finalized state of the system after
// the Nucleus has completed its initialization.
type BootSequence struct {
	PlatformID platforms.PlatformClass `json:"platform_id"` // e.g., Vehicle, Laptop
	TrustScore float64                 `json:"trust_score"` // 0.0 to 1.0 (Bayesian)
	IsVerified bool                    `json:"is_verified"` // Attestation result
	Mode       string                  `json:"mode"`        // Autonomous | Discovery | Safe
	UserRole   string                  `json:"user_role"`   // Operator | Admin
}


type EnvConfig struct {
	SchemaVersion int                `json:"schema_version"`
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
    Type    string  `json:"type"`    // CPU, GPU, TPU
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

// PlatformScore tracks the heuristic weight for a specific platform type.
type PlatformScore struct {
	Class      PlatformClass `json:"class"`
	Score      float64       `json:"score"`      // Raw cumulative score
	MaxScore   float64       `json:"max_score"`  // Potential maximum for normalization
	Confidence uint16        `json:"confidence"` // Normalized Q16 (0-65535)
	Signals    []string      `json:"signals"`    // Evidence found (e.g., "CAN_BUS_PRESENT")
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







