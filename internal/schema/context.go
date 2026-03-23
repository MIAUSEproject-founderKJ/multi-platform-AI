//internal\schema\context.go

package schema

import "go.uber.org/zap"

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
