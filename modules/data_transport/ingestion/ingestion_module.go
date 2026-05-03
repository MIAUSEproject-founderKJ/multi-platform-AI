//modules/data_transport/ingestion/ingestion_module.go

package transport_ingestion

import (
	"context"
	"sync/atomic"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/file"
	kernel_lifecycle "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/lifecycle"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

type IngestionModule struct {
	kernel_lifecycle.BaseModule
	repo    file.FileRepository
	bootctx runtime_types.ExecutionContext

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

func (m *IngestionModule) Allowed(bootctx runtime_types.ExecutionContext) bool { return true }

func (m *IngestionModule) Init(bootctx runtime_types.ExecutionContext) error { return nil }

func (m *IngestionModule) Start() error { return nil }

func (m *IngestionModule) Stop() error { return nil }

func (m *IngestionModule) Run(bootctx context.Context) error {
	<-bootctx.Done()
	return nil
}

func (m *IngestionModule) SupportedPlatforms() []internal_environment.PlatformClass { return nil }

func (m *IngestionModule) RequiredCapabilities() internal_environment.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}

func (m *IngestionModule) Optional() bool { return false }

func (m *IngestionModule) handle(bootctx context.Context, payload []byte) error {
	return m.repo.StoreChunk(bootctx, payload)
}
