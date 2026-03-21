//modules/audit_module.go

package modules

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

type AuditModule struct {
	ctx *schema.RuntimeContext
}

func (m *AuditModule) Allowed(ctx *schema.BootContext) bool {
	return true
}

func (m *AuditModule) SetRuntime(rtx *schema.RuntimeContext) {
	m.ctx = rtx
}
func NewAuditModule() DomainModule {
	return &AuditModule{}
}
