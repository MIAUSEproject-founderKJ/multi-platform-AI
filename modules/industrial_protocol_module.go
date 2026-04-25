//modules/industrial_protocol_module.go

package modules

import (
	"context"
	"sync/atomic"

	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

type IndustrialProtocolModule struct {
	BaseModule // embed BaseModule
	healthy    atomic.Bool
	ctx        *schema_boot.BootContext // store BootContext locally
}

// Factory function
func NewIndustrialProtocolModule() DomainModule {
	return &IndustrialProtocolModule{
		BaseModule: BaseModule{
			name: "IndustrialProtocolModule",
			deps: []string{"TelemetryModule"},
		},
	}
}

func (m *IndustrialProtocolModule) Init(ctx *schema_boot.BootContext) error {
	m.ctx = ctx // store for internal use
	m.healthy.Store(true)
	return nil
}

func (m *IndustrialProtocolModule) Name() string {
	return "IndustrialProtocolModule"
}

func (m *IndustrialProtocolModule) Category() ModuleCategory {
	return ModuleDomain
}

func (m *IndustrialProtocolModule) DependsOn() []string {
	return []string{"TelemetryModule"}
}

func (m *IndustrialProtocolModule) Allowed(ctx *schema_boot.BootContext) bool {
	return true
}

func (m *IndustrialProtocolModule) Start() error {
	return nil
}

func (m *IndustrialProtocolModule) Stop() error {
	return nil
}

func (m *IndustrialProtocolModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *IndustrialProtocolModule) Handle(ctx context.Context, payload []byte) error {
	if m.ctx != nil && m.ctx.Logger != nil {
		m.ctx.Logger.Info("industrial protocol message")
	}
	return nil
}

func (m *IndustrialProtocolModule) Healthy() bool {
	return m.healthy.Load()
}

func (m *IndustrialProtocolModule) SupportedPlatforms() []schema_system.PlatformClass { return nil }
func (m *IndustrialProtocolModule) RequiredCapabilities() schema_security.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}
func (m *IndustrialProtocolModule) Optional() bool { return false }
