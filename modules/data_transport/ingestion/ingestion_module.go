//modules/data_transport/ingestion/ingestion_module.go

package transport_ingestion

import (
	"context"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/file"
	kernel_lifecycle "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/lifecycle"
)

type IngestionModule struct {
	kernel_lifecycle.BaseModule
	repo file.FileRepository
	ctx  *bootstrap.BootContext

	running atomic.Bool
}

func NewIngestionModule() domain_shared.DomainModule {
	return &IngestionModule{
		BaseModule: kernel_lifecycle.BaseModule{
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

func (m *IngestionModule) Allowed(*bootstrap.BootContext) bool { return true }

func (m *IngestionModule) Init(*bootstrap.BootContext) error { return nil }

func (m *IngestionModule) Start() error { return nil }

func (m *IngestionModule) Stop() error { return nil }

func (m *IngestionModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *IngestionModule) SupportedPlatforms() []internal_environment.PlatformClass { return nil }

func (m *IngestionModule) RequiredCapabilities() internal_environment.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}

func (m *IngestionModule) Optional() bool { return false }

func (m *IngestionModule) handle(ctx context.Context, payload []byte) error {
	return m.repo.StoreChunk(ctx, payload)
}
