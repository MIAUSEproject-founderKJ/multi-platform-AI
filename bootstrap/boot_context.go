//MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/boot_context.go

package bootstrap

import (
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
	"go.uber.org/zap"
)

type BootContext struct {
	// ===== Identity / Classification =====
	platformClass internal_environment.PlatformClass
	service       user_setting.ServiceType
	entity        internal_environment.EntityKind
	tier          user_setting.TierType
	bootMode      internal_boot.BootMode

	// ===== Security =====
	trustLevel  user_setting.TrustLevel
	permissions map[user_setting.PermissionKey]bool
	permMask    internal_verification.PermissionMask

	// ===== Capability =====
	capabilities internal_environment.CapabilitySet

	// ===== Infrastructure (INTERNAL ONLY) =====
	vault  verification_persistence.VaultStore
	logger *zap.Logger
}
