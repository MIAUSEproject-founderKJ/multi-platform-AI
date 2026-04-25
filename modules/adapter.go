//modules/adapter.go
/*
Keeps your existing modules unchanged
Bridges old interface → new runtime
Allows system to compile and run immediately
You can migrate modules incrementally later
*/
package modules

import (
	"context"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/supervisor"
)

// Adapter wraps old DomainModule into new supervisor.Module
type Adapter struct {
	legacy DomainModule
}

// AdaptModules converts legacy modules into supervisor.Module
func AdaptModules(mods []DomainModule, rtx *engine.RuntimeContext) []supervisor.Module {
	out := make([]supervisor.Module, 0, len(mods))

	for _, m := range mods {
		// Inject runtime if supported
		if rm, ok := m.(RuntimeAware); ok {
			rm.SetRuntime(rtx)
		}

		out = append(out, &Adapter{legacy: m})
	}

	return out
}

func (a *Adapter) Name() string {
	return a.legacy.Name()
}

func (a *Adapter) Init(ctx context.Context) error {
	// Temporary: BootContext not passed yet
	return a.legacy.Init(nil)
}

func (a *Adapter) Start(ctx context.Context) error {
	// Map Start -> Run
	return a.legacy.Run(ctx)
}

func (a *Adapter) Stop(ctx context.Context) error {
	// No-op for now (legacy modules do not support shutdown)
	return nil
}

func (a *Adapter) Health() error {
	// Assume healthy unless Run fails
	return nil
}
