//modules/ingestion_module.go

package modules

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
)

type IngestionModule struct {
	ctx     *runtime.ExecutionContext
	healthy atomic.Bool
}

func NewIngestionModule() DomainModule {
	return &IngestionModule{}
}

func (m *IngestionModule) Name() string {
	return "IngestionModule"
}

func (m *IngestionModule) DependsOn() []string {
	return nil
}

func (m *IngestionModule) Init(ctx *runtime.ExecutionContext) error {
	m.ctx = ctx
	m.healthy.Store(true)
	return nil
}

func (m *IngestionModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *IngestionModule) Handle(ctx context.Context, payload []byte) error {

	if len(payload) == 0 {
		return errors.New("empty ingestion payload")
	}

	var raw map[string]interface{}

	if err := json.Unmarshal(payload, &raw); err != nil {
		return err
	}

	// forward to telemetry pipeline
	return m.ctx.Router.Publish("telemetry", payload)
}

func (m *IngestionModule) Healthy() bool {
	return m.healthy.Load()
}