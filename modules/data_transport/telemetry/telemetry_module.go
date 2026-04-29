// modules/data_transport/telemetry/telemetry_module.go exports metrics to network.

package transport_telemetry

import (
	"context"
	"sync/atomic"

	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
	kernel_lifecycle "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/lifecycle"
)

type TelemetryModule struct {
	BaseModule

	ctx     *internal_boot.BootContext
	client  TelemetryClient
	running atomic.Bool
	healthy atomic.Bool
}

type TelemetryClient interface {
	Send([]byte) error
}

// Factory
func NewTelemetryModule() domain_shared.DomainModule {
	return &TelemetryModule{
		BaseModule: kernel_lifecycle.BaseModule{
			name: "TelemetryModule",
			deps: []string{"IngestionModule"},
		},
	}
}

func (m *TelemetryModule) Init(ctx *internal_boot.BootContext) error {
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
func (m *TelemetryModule) DependsOn() []string                         { return []string{"IngestionModule"} }
func (m *TelemetryModule) Allowed(ctx *internal_boot.BootContext) bool { return true }
func (m *TelemetryModule) Start() error                                { return nil }
func (m *TelemetryModule) Stop() error                                 { return nil }
func (m *TelemetryModule) Healthy() bool                               { return m.healthy.Load() }
func (m *TelemetryModule) SupportedPlatforms() []internal_environment.PlatformClass {
	return nil
}

func (m *TelemetryModule) RequiredCapabilities() internal_environment.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}
func (m *TelemetryModule) Optional() bool { return false }
