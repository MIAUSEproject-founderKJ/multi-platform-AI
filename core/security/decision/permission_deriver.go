// core/security/decision/permission_deriver.go
package security_decision

import (
	bootstrap_resolver "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/resolver"
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type PermissionResolver interface {
	Resolve(ctx *AuthorizationContext) []user_setting.PermissionKey
}

type AuthorizationContext struct {
	Platform internal_common.PlatformClass
	Entity   internal_common.EntityKind
	Tier     internal_common.TierType
	Service  user_setting.ServiceType
}

type DefaultPermissionResolver struct{}

func (r *DefaultPermissionResolver) Resolve(ctx *AuthorizationContext) []user_setting.PermissionKey {
	return bootstrap_resolver.DerivePermissions(
		ctx.Platform,
		ctx.Entity,
		ctx.Tier,
		ctx.Service,
	)
}
