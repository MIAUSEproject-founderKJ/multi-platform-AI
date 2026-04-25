//internal/schema/boot/context.go

package schema_boot

import (
	security_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
	"go.uber.org/zap"
)

type BootContext struct {
	PlatformClass schema_system.PlatformClass
	Capabilities  schema_security.CapabilitySet
	Vault         security_persistence.VaultStore
	Service       schema_identity.ServiceType
	Entity        schema_system.EntityType
	Tier          schema_identity.TierType
	BootMode      BootMode
	Logger        *zap.Logger
	Permissions   map[schema_identity.Permission]bool // storage
	PermMask      schema_security.PermissionMask      // runtime
	TrustLevel    schema_identity.TrustLevel
}
