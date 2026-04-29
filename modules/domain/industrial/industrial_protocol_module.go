// modules/domain/industrial/industrial_protocol_module.go
package module_industrial

import (
	"context"
	"sync/atomic"

	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
	kernel_lifecycle "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/lifecycle"
)

type IndustrialProtocolModule struct {
	BaseModule // embed BaseModule
	healthy    atomic.Bool
	ctx        *bootstrap.BootContext // store BootContext locally
}

// Factory function
func NewIndustrialProtocolModule() domain_shared.DomainModule {
	return &IndustrialProtocolModule{
		BaseModule: kernel_lifecycle.BaseModule{
			name: "IndustrialProtocolModule",
			deps: []string{"TelemetryModule"},
		},
	}
}

func (m *IndustrialProtocolModule) Init(ctx *bootstrap.BootContext) error {
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

func (m *IndustrialProtocolModule) Allowed(ctx *bootstrap.BootContext) bool {
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

func (m *IndustrialProtocolModule) SupportedPlatforms() []internal_environment.PlatformClass {
	return nil
}
func (m *IndustrialProtocolModule) RequiredCapabilities() internal_environment.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}
func (m *IndustrialProtocolModule) Optional() bool { return false }
