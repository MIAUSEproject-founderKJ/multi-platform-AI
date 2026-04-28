// core/verification/decision/permission_deriver.go
package verification_decision

import (
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type PermissionResolver interface {
	Resolve(ctx *AuthorizationContext) []user_setting.PermissionKey
}

type AuthorizationContext struct {
	Platform internal_environment.PlatformClass
	Entity   internal_environment.EntityKind
	Tier     user_setting.TierType
	Service  user_setting.ServiceType
}

type DefaultPermissionResolver struct{}

func (r *DefaultPermissionResolver) Resolve(ctx *AuthorizationContext) []user_setting.PermissionKey {
	return DerivePermissions(
		ctx.Platform,
		ctx.Entity,
		ctx.Tier,
		ctx.Service,
	)
}
