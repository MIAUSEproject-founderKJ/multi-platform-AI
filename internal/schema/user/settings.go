//internal/schema/user/settings.go

package user_setting

import (
	"fmt"
	"time"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
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

func (s *UserSession) HasPermission(p internal_verification.PermissionMask) bool {
	if s == nil {
		return false
	}
	return s.PermMask&p != 0
}

type UserCoreConfig struct {
	MainLang      string
	PowerMode     string
	PrivacyMode   string
	UpdateMode    string
	PreferredMode string
}

type CustomizedConfig struct {
	Version      string
	LastModified time.Time
	UserCoreConfig
}

type UserIdentity struct {
	Username string
}

type UserPreferences struct {
	AIStyle         string
	AutoSave        bool
	EnableTelemetry bool
}

type UserSession struct {
	Identity    *UserIdentity
	Config      UserCoreConfig
	Preferences *UserPreferences
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

type PermissionKey string

const (
	PermUser  PermissionKey = "user"
	PermAdmin PermissionKey = "admin"

	PermBasicRuntime   PermissionKey = "basic_runtime"
	PermConfigEdit     PermissionKey = "config_edit"
	PermDiagnostics    PermissionKey = "diagnostics"
	PermHardwareIO     PermissionKey = "hardware_io"
	PermSafetyOverride PermissionKey = "safety_override"
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
	permissions map[PermissionKey]bool,
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
	Platform internal_environment.PlatformClass
	Entity   internal_environment.EntityKind
	Tier     TierType
	Service  ServiceType
}
