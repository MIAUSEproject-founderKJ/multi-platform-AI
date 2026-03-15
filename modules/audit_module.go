//modules/audit_module.go

package modules

type AuditModule struct {
}

func NewAuditModule() DomainModule {
	return &AuditModule{}
}
