//modules/industrial_protocol_module.go

package modules

import (
	"context"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type IndustrialProtocolModule struct {
	BaseModule // embed BaseModule
	healthy    atomic.Bool
	ctx        *schema.BootContext // store BootContext locally
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

func (m *IndustrialProtocolModule) Init(ctx *schema.BootContext) error {
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

func (m *IndustrialProtocolModule) Allowed(ctx *schema.BootContext) bool {
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

func (m *IndustrialProtocolModule) SupportedPlatforms() []schema.PlatformClass { return nil }
func (m *IndustrialProtocolModule) RequiredCapabilities() schema.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}
func (m *IndustrialProtocolModule) Optional() bool { return false }
