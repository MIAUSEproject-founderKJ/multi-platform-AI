//core\security\decision\decision_engine.go

package security

import (
	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
)

type Enforcer struct {
	ctx *schema_boot.BootContext
}

func NewEnforcer(ctx *schema_boot.BootContext) *Enforcer {
	return &Enforcer{ctx: ctx}
}

func (e *Enforcer) Allow(p schema_identity.Permission) bool {

	v, ok := e.ctx.Permissions[p]
	if !ok {
		return false
	}

	return v
}
