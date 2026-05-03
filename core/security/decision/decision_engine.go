//core/verification/decision/decision_engine.go

package security_decision

import (
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

type Enforcer struct {
	ctx runtime_types.ExecutionContext
}

func NewEnforcer(bootctx runtime_types.ExecutionContext) *Enforcer {
	return &Enforcer{ctx: bootctx}
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

func (as *AuthorizationService) Authorize(authbootctx *AuthorizationContext) map[user_setting.PermissionKey]bool {
	perms := as.Resolver.Resolve(authbootctx)

	permMap := make(map[user_setting.PermissionKey]bool)
	for _, p := range perms {
		permMap[p] = true
	}

	return permMap
}
