//core/verification/decision/decision_engine.go

package verification_decision

import (
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type Enforcer struct {
	ctx *internal_boot.BootContext
}

func NewEnforcer(ctx *internal_boot.BootContext) *Enforcer {
	return &Enforcer{ctx: ctx}
}

func (e *Enforcer) Allow(p user_setting.PermissionKey) bool {

	v, ok := e.ctx.Permissions[p]
	if !ok {
		return false
	}

	return v
}

type AuthorizationService struct {
	Resolver PermissionResolver
}

func (as *AuthorizationService) Authorize(authCtx *AuthorizationContext) map[user_setting.PermissionKey]bool {
	perms := as.Resolver.Resolve(authCtx)

	permMap := make(map[user_setting.PermissionKey]bool)
	for _, p := range perms {
		permMap[p] = true
	}

	return permMap
}
