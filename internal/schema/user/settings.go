//internal/schema/user/settings.go

package user_setting

import (
	"fmt"
	"time"

	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
)

type InteractionCapability struct {
	CLI   bool
	TUI   bool
	GUI   bool
	Voice bool
}

// Move these from interaction to schema so everyone can see them
type InteractionMode string

const (
	ModeCLIonly InteractionMode = "cli"
	ModeTUIonly InteractionMode = "tui"
	ModeGUIonly InteractionMode = "gui"
	ModeFull    InteractionMode = "cgtvio" // Console + GUI + TUI + Voice
	ModeCT      InteractionMode = "ct"     // Console + TUI
	ModeCG      InteractionMode = "cg"     // Console + GUI
	ModeCVo     InteractionMode = "cvo"    // Console + VoiceOut
	ModeCVi     InteractionMode = "cvi"    // Console + VoiceInOut
	ModeTG      InteractionMode = "tg"     // TUI + GUI
	ModeTVi     InteractionMode = "tvi"    // TUI + VoiceIn
	ModeTVo     InteractionMode = "tvo"    // TUI + VoiceOut
	ModeGVo     InteractionMode = "gvo"    // GUI + VoiceOut
	ModeGVi     InteractionMode = "gvi"    // GUI + VoiceIn
	ModeGVio    InteractionMode = "gvio"   // GUI + VoiceIn+Out
	ModeTVio    InteractionMode = "tvio"   // TUI + VoiceInOut
	ModeGTVio   InteractionMode = "gtvio"  // GUI + TUI + VoiceInOut
	ModeVio     InteractionMode = "vio"    // VoiceInOut only
)

// Define what an Orchestrator DOES, not what it IS.
type Orchestrator interface {
	StartAll(session *UserSession)
	Broadcast(msg string)
}

func (s *UserSession) HasPermission(p PermissionKey) bool {
	if s == nil {
		return false
	}
	return s.Claims.Permissions[p]
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

	Claims SessionClaims
}

type SessionClaims struct {
	SessionID string

	Platform internal_common.PlatformClass
	Entity   internal_common.EntityKind

	Tier    internal_common.TierType
	Service ServiceType

	Permissions map[PermissionKey]bool

	CreatedAt time.Time
	ExpiresAt time.Time
}

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
) *SessionClaims {

	return &SessionClaims{
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
	Platform internal_common.PlatformClass
	Entity   internal_common.EntityKind
	Tier     internal_common.TierType
	Service  ServiceType
}
