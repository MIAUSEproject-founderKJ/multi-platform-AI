//internal/schema/context.go

package schema

type BootContext struct {
	PlatformClass PlatformClass
	Capabilities  CapabilitySet

	Service  ServiceType
	Entity   EntityType
	Tier     TierType
	BootMode BootMode
	Logger
	Permissions map[Permission]bool
	TrustLevel  TrustLevel
}
