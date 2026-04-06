//internal/schema/user.go

package schema

import "time"

type UserSession struct {
	SessionID string

	Platform PlatformClass
	Entity   EntityType
	Tier     TierType
	Service  ServiceType

	Permissions map[Permission]bool

	Config       *CustomizedConfig
	Capabilities CapabilitySet
	Mode         string

	CreatedAt time.Time
	ExpiresAt time.Time
	CapProfile    *CapabilityProfile
	Mode          string

	Orchestrator interface{} // runtime binding
}

// ------------------------------------------------------------
// Tier System
// ------------------------------------------------------------
//Use TierType (string) externally for readability and compatibility. Use EntityType (uint8) internally for speed and clarity.
type TierType string

const (
	TierUnknown    TierType = "unknown"
	TierPersonal   TierType = "personal"
	TierEnterprise TierType = "enterprise"
	TierTester     TierType = "tester"
)

// Optional richer structure
type TierProfile struct {
	Name TierType
}

// ------------------------------------------------------------
// Service System
// ------------------------------------------------------------

type ServiceType string

const (
	ServiceUnknown    ServiceType = "unknown"
	ServicePersonal   ServiceType = "personal_ai"
	ServiceEnterprise ServiceType = "enterprise_ai"
	ServiceSystem     ServiceType = "system_runtime"
	ServiceIndustrial ServiceType = "industrial_control"
	ServiceMobility   ServiceType = "autonomous_mobility"
)

// Optional richer structure
type ServiceProfile struct {
	Name ServiceType
}

type Attestation struct {
	SessionToken string
	Valid        bool
	Level        TrustLevel
}

type TrustLevel uint8

const (
	TrustUntrusted TrustLevel = iota
	TrustUser
	TrustDevice
	TrustAdmin
	TrustSystem
)

type Permission string

const (
	PermUser  Permission = "user"
	PermAdmin Permission = "admin"

	PermBasicRuntime   Permission = "basic_runtime"
	PermConfigEdit     Permission = "config_edit"
	PermDiagnostics    Permission = "diagnostics"
	PermHardwareIO     Permission = "hardware_io"
	PermSafetyOverride Permission = "safety_override"
)

type UserConfig struct {
	Username        string
	PreferredMode   string // "cli", "voice", "hybrid"
	AIStyle         string // "concise", "balanced", "verbose"
	AutoSave        bool
	EnableTelemetry bool

	Runtime CustomizedConfig
}
