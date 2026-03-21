// modules/telemetry_module.go exports metrics to network.
package modules

import (
	"context"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type TelemetryModule struct {
	BaseModule

	ctx *schema.BootContext

	client  TelemetryClient
	running atomic.Bool
}

type TelemetryClient interface {
	Send([]byte) error
}

func NewTelemetryModule() DomainModule {

	m := &TelemetryModule{
		BaseModule: BaseModule{
			name: "TelemetryModule",
			deps: []string{"IngestionModule"},
		},
	}

	return m
}

func (m *TelemetryModule) Init(ctx *schema.BootContext) error {

	m.ctx = ctx

	m.setHealthy(true)

	ctx.Logger.Info("TelemetryModule initialized")

	return nil
}

func (m *TelemetryModule) Run(ctx context.Context) error {

	m.running.Store(true)

	m.ctx.Logger.Info("TelemetryModule started")

	<-ctx.Done()

	m.running.Store(false)

	m.ctx.Logger.Info("TelemetryModule stopped")

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
