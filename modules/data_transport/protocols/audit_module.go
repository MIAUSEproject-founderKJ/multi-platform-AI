// modules/data_transport/protocols/audit_module.go

package transport_audit

import (
	"context"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
	kernel_lifecycle "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/lifecycle"
	runtime_engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
)

type AuditModule struct {
	kernel_lifecycle.BaseModule
	runtime *runtime_engine.RuntimeContext
	cancel  context.CancelFunc
}

func NewAuditModule() domain_shared.DomainModule {
	m := &AuditModule{}
	m.SetName("AuditModule")
	return m
}

func (m *AuditModule) Start() error {
	// No-op start for now
	return nil
}

func (m *AuditModule) RequiredCapabilities() internal_environment.CapabilitySet {
	return 0 // no requirement → always eligible
}

func (m *AuditModule) Optional() bool {
	return true
}

func (m *AuditModule) Allowed(ctx *bootstrap.BootContext) bool {
	return ctx.Permissions[user_setting.PermDiagnostics]
}

func (m *AuditModule) Category() ModuleCategory {
	return ModuleDomain
}

func (m *AuditModule) SupportedPlatforms() []internal_environment.PlatformClass {
	return nil // capability-driven only
}

func (m *AuditModule) DependsOn() []string {
	return nil
}

func (m *AuditModule) Init(ctx *bootstrap.BootContext) error {
	if m.runtime == nil {
		return fmt.Errorf("runtime not set")
	}
	return nil
}

func (m *AuditModule) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	ch := m.runtime.Bus.Subscribe("audit.events")

	for {
		select {
		case <-ctx.Done():
			return nil

		case msg := <-ch:
			fmt.Printf("[Audit] Topic=%s Data=%s\n", msg.Topic, string(msg.Data))
		}
	}
}

func (m *AuditModule) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

func (m *AuditModule) SetRuntime(rtx *runtime_engine.RuntimeContext) {
	m.runtime = rtx
}
