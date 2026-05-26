//MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/boot_context.go

package bootstrap

import (
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	"go.uber.org/zap"
)

// boot context is the shared state during the boot process, passed through orchestrator and resolvers. It contains all the information needed to make decisions and build the final execution context. It is NOT the final execution context, but a mutable state that evolves during boot.
type BootContext struct {
	platformClass internal_environment.PlatformClass
	service       user_setting.ServiceType
	entity        internal_environment.EntityKind
	tier          user_setting.TierType
	bootMode      internal_boot.BootMode

	trustLevel   user_setting.TrustLevel
	capabilities internal_environment.CapabilitySet

	vault  verification_persistence.VaultStore
	logger *zap.Logger
}
