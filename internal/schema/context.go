//internal\schema\context.go

package schema

import (
	"database/sql"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	"go.uber.org/zap"
)

type BootContext struct {
	PlatformClass PlatformClass
	Capabilities  CapabilitySet

	Service  ServiceType
	Entity   EntityType
	Tier     TierType
	BootMode BootMode

	Permissions map[Permission]bool
	TrustLevel  TrustLevel

	Logger *zap.Logger
}

type RuntimeContext struct {
	Router *router.Router
	Bus    *MessageBus
	DB     *sql.DB
	Logger *zap.Logger
}
