//internal/schema/user.go

package schema

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/interaction"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func (s *UserSession) HasPermission(p PermissionMask) bool {
	if s == nil {
		return false
	}
	return s.PermMask&p != 0
}

type UserSession struct {
	SessionID string

	Platform PlatformClass
	Entity   EntityType
	Tier     TierType
	Service  ServiceType

	Permissions map[Permission]bool // storage
	PermMask    PermissionMask      // runtime

	Config *UserConfig

	Mode       interaction.InteractionMode
	CapProfile *CapabilityProfile

	Attestation *Attestation

	CreatedAt time.Time
	ExpiresAt time.Time

	Capabilities schema.CapabilitySet
	Orchestrator interaction.Orchestrator
}

// ------------------------------------------------------------
// Tier System
// ------------------------------------------------------------
// Use TierType (string) externally for readability and compatibility. Use EntityType (uint8) internally for speed and clarity.
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
	PreferredMode   string
	AIStyle         string
	AutoSave        bool
	EnableTelemetry bool

	// Merge runtime config directly
	MainLang    string
	PowerMode   string
	PrivacyMode string
	UpdateMode  string
}
