//modules/adapter/module_adapter.go
/*
Keeps your existing modules unchanged
Bridges old interface → new runtime
Allows system to compile and run immediately
You can migrate modules incrementally later
*/
package modules_adapter

import (
	"context"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/supervisor"
)

// Adapter wraps old DomainModule into new supervisor.Module
type Adapter struct {
	legacy shared.DomainModule
}

func AdaptModules(mods []shared.DomainModule, rtx *engine.RuntimeContext) []supervisor.Module {
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
	return a.legacy.Init(ctx)
}

func (a *Adapter) Start(ctx context.Context) error {
	// Map Start -> Run
	return a.legacy.Run(ctx)
}

type Stoppable interface {
	Stop(ctx context.Context) error
}

func (a *Adapter) Stop(ctx context.Context) error {
	if s, ok := a.legacy.(Stoppable); ok {
		return s.Stop(ctx)
	}
	return nil
}

type HealthAware interface {
	Health() error
}

func (a *Adapter) Health() error {
	if h, ok := a.legacy.(HealthAware); ok {
		return h.Health()
	}
	return nil
}
