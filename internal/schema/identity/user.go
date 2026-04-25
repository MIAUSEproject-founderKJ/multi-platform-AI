//internal/schema/identity/user.go

package schema_identity

import (
	"time"

	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

// Move these from interaction to schema so everyone can see them
type InteractionMode string

const (
	ModeCLI   InteractionMode = "cli"
	ModeTUI   InteractionMode = "tui"
	ModeGUI   InteractionMode = "gui"
	ModeVoice InteractionMode = "voice"
)

// Define what an Orchestrator DOES, not what it IS.
type Orchestrator interface {
	StartAll(session *UserSession)
	Broadcast(msg string)
}

func (s *UserSession) HasPermission(p schema_security.PermissionMask) bool {
	if s == nil {
		return false
	}
	return s.PermMask&p != 0
}

type UserSession struct {
	SessionID string

	Platform schema_system.PlatformClass
	Entity   schema_system.EntityType
	Tier     TierType
	Service  ServiceType

	Permissions map[Permission]bool            // storage
	PermMask    schema_security.PermissionMask // runtime

	Config *UserConfig

	Mode       InteractionMode
	CapProfile *schema_security.CapabilityProfile

	Attestation *Attestation

	CreatedAt time.Time
	ExpiresAt time.Time

	Capabilities schema_security.CapabilitySet
	Orchestrator Orchestrator
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


type SessionBuilder struct{}

func (b *SessionBuilder) Build(
    ctx *BuildContext,
    permissions map[Permission]bool,
) *UserSession {

    return &UserSession{
        SessionID:   fmt.Sprintf("%d", time.Now().UnixNano()),
        Platform:    ctx.Platform,
        Entity:      ctx.Entity,
        Tier:        ctx.Tier,
        Service:     ctx.Service,
        Permissions: permissions,
        CreatedAt:   time.Now(),
        ExpiresAt:   time.Now().Add(24 * time.Hour),
    }
}

type BuildContext struct {
    Platform schema_system.PlatformClass
    Entity   schema_system.EntityType
    Tier     schema_identity.TierType
    Service  schema_identity.ServiceType
}