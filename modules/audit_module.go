// modules/audit_module.go

package modules

import (
	"context"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime"
)

type AuditModule struct {
	BaseModule
	runtime *runtime.RuntimeContext
	cancel  context.CancelFunc
}

func NewAuditModule() DomainModule {
	m := &AuditModule{}
	m.SetName("AuditModule")
	return m
}

func (m *AuditModule) Start() error {
	// No-op start for now
	return nil
}

func (m *AuditModule) RequiredCapabilities() schema.CapabilitySet {
	return 0 // no requirement → always eligible
}

func (m *AuditModule) Optional() bool {
	return true
}

func (m *AuditModule) Allowed(ctx *schema.BootContext) bool {
	return ctx.Permissions[schema.PermDiagnostics]
}

func (m *AuditModule) Category() ModuleCategory {
	return ModuleDomain
}

func (m *AuditModule) SupportedPlatforms() []schema.PlatformClass {
	return nil // capability-driven only
}

func (m *AuditModule) DependsOn() []string {
	return nil
}

func (m *AuditModule) Init(ctx *schema.BootContext) error {
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

func (m *AuditModule) SetRuntime(rtx *runtime.RuntimeContext) {
	m.runtime = rtx
}
