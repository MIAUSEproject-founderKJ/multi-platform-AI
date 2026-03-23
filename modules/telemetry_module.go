// modules/telemetry_module.go exports metrics to network.

package modules

import (
	"context"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type TelemetryModule struct {
	BaseModule

	ctx     *schema.BootContext
	client  TelemetryClient
	running atomic.Bool
	healthy atomic.Bool
}

type TelemetryClient interface {
	Send([]byte) error
}

// Factory
func NewTelemetryModule() DomainModule {
	return &TelemetryModule{
		BaseModule: BaseModule{
			name: "TelemetryModule",
			deps: []string{"IngestionModule"},
		},
	}
}

func (m *TelemetryModule) Init(ctx *schema.BootContext) error {
	m.ctx = ctx
	m.healthy.Store(true)

	if m.ctx != nil && m.ctx.Logger != nil {
		m.ctx.Logger.Info("TelemetryModule initialized")
	}

	return nil
}

func (m *TelemetryModule) Run(ctx context.Context) error {
	m.running.Store(true)

	if m.ctx != nil && m.ctx.Logger != nil {
		m.ctx.Logger.Info("TelemetryModule started")
	}

	<-ctx.Done()

	m.running.Store(false)

	if m.ctx != nil && m.ctx.Logger != nil {
		m.ctx.Logger.Info("TelemetryModule stopped")
	}

	return nil
}

func (m *TelemetryModule) Handle(ctx context.Context, payload []byte) error {
	if len(payload) == 0 {
		return nil
	}

	if m.client != nil {
		return m.client.Send(payload)
	}

	return nil
}

// DomainModule compliance
func (m *TelemetryModule) Name() string { return "TelemetryModule" }
func (m *TelemetryModule) Category() ModuleCategory {
	return ModuleDomain
}
func (m *TelemetryModule) DependsOn() []string                  { return []string{"IngestionModule"} }
func (m *TelemetryModule) Allowed(ctx *schema.BootContext) bool { return true }
func (m *TelemetryModule) Start() error                         { return nil }
func (m *TelemetryModule) Stop() error                          { return nil }
func (m *TelemetryModule) Healthy() bool                        { return m.healthy.Load() }
func (m *TelemetryModule) SupportedPlatforms() []schema.PlatformClass {
	return nil
}

func (m *TelemetryModule) RequiredCapabilities() schema.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}
func (m *TelemetryModule) Optional() bool { return false }
