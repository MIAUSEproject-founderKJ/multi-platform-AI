//bootstrap/boot_context.go

package bootstrap

import (
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/verification/persistence"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
	"go.uber.org/zap"
)

type BootContext struct {
	PlatformClass internal_environment.PlatformClass
	Capabilities  internal_environment.CapabilitySet
	Vault         verification_persistence.VaultStore
	Service       user_setting.ServiceType
	Entity        internal_environment.EntityKind
	Tier          user_setting.TierType
	BootMode      internal_boot.BootMode
	Logger        *zap.Logger
	Permissions   map[user_setting.PermissionKey]bool  // storage
	PermMask      internal_verification.PermissionMask // runtime
	TrustLevel    user_setting.TrustLevel
}
