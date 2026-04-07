//internal/schema/context.go

package schema

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"

	"go.uber.org/zap"
)

type BootContext struct {
	PlatformClass PlatformClass
	Capabilities  CapabilitySet
	Vault         security.VaultStore
	Service       ServiceType
	Entity        EntityType
	Tier          TierType
	BootMode      BootMode
	Logger        *zap.Logger
	Permissions   map[Permission]bool
	TrustLevel    TrustLevel
}
