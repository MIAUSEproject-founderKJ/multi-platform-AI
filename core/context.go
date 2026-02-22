//core/context.go

package core

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/schema"

type RuntimeContext struct {
	PlatformClass schema.PlatformClass
	Capabilities  CapabilitySet
	Service       ServiceType
	Entity        EntityType
	Tier          TierType
	BootMode      schema.BootMode
	Permissions   PermissionSet
}