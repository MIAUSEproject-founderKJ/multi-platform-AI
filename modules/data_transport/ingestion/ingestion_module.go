//modules/data_transport/ingestion/ingestion_module.go

package transport_ingestion

import (
	"context"
	"sync/atomic"

	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/file"
)

type IngestionModule struct {
	BaseModule
	repo file.FileRepository
	ctx  *internal_boot.BootContext

	running atomic.Bool
}

func NewIngestionModule() DomainModule {
	return &IngestionModule{
		BaseModule: BaseModule{
			name: "IngestionModule",
			deps: []string{},
		},
	}
}

func (m *IngestionModule) Name() string { return m.name }

func (m *IngestionModule) Category() ModuleCategory {
	return ModuleDomain
}

func (m *IngestionModule) DependsOn() []string { return m.deps }

func (m *IngestionModule) Allowed(*internal_boot.BootContext) bool { return true }

func (m *IngestionModule) Init(*internal_boot.BootContext) error { return nil }

func (m *IngestionModule) Start() error { return nil }

func (m *IngestionModule) Stop() error { return nil }

func (m *IngestionModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *IngestionModule) SupportedPlatforms() []internal_environment.PlatformClass { return nil }

func (m *IngestionModule) RequiredCapabilities() internal_verification.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}

func (m *IngestionModule) Optional() bool { return false }

func (m *IngestionModule) handle(ctx context.Context, payload []byte) error {
	return m.repo.StoreChunk(ctx, payload)
}
