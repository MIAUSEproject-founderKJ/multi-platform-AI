//core/security/decision/permission_deriver.go
package security_decision

import schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
import schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"

type PermissionResolver interface {
    Resolve(ctx *AuthorizationContext) []schema_identity.Permission
}

type AuthorizationContext struct {
    Platform schema_system.PlatformClass
    Entity   schema_system.EntityType
    Tier     schema_identity.TierType
    Service  schema_identity.ServiceType
}

type DefaultPermissionResolver struct{}

func (r *DefaultPermissionResolver) Resolve(ctx *AuthorizationContext) []schema_identity.Permission {
    return DerivePermissions(
        ctx.Platform,
        ctx.Entity,
        ctx.Tier,
        ctx.Service,
    )
}

