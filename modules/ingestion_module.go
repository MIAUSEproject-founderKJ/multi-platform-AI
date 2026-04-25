//modules/ingestion_module.go

package modules

import (
	"context"
	"sync/atomic"

	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/file"
)

type IngestionModule struct {
	BaseModule
	repo file.FileRepository
	ctx  *schema_boot.BootContext

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

func (m *IngestionModule) Allowed(*schema_boot.BootContext) bool { return true }

func (m *IngestionModule) Init(*schema_boot.BootContext) error { return nil }

func (m *IngestionModule) Start() error { return nil }

func (m *IngestionModule) Stop() error { return nil }

func (m *IngestionModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *IngestionModule) SupportedPlatforms() []schema_system.PlatformClass { return nil }

func (m *IngestionModule) RequiredCapabilities() schema_security.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}

func (m *IngestionModule) Optional() bool { return false }

func (m *IngestionModule) handle(ctx context.Context, payload []byte) error {
	return m.repo.StoreChunk(ctx, payload)
}
