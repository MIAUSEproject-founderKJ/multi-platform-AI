//core\security\policy_enforcer.go

package security

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

type Enforcer struct {
	ctx *schema.BootContext
}

func NewEnforcer(ctx *schema.BootContext) *Enforcer {
	return &Enforcer{ctx: ctx}
}

func (e *Enforcer) Allow(p schema.Permission) bool {

	v, ok := e.ctx.Permissions[p]
	if !ok {
		return false
	}

	return v
}
