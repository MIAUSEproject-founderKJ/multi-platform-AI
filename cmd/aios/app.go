//cmd/aios/app.go

package main

import (
	"context"
	"errors"
	"net/http"

	modules_adapter "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/adapter"
	kernel_registry "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/registry"
	kernel_supervisor "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/supervisor"
	runtime_engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
	runtime_supervisor "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/supervisor"
	"go.uber.org/zap"
)

// ============================================================
// APP COMPOSITION
// ============================================================

type App struct {
	log        *zap.Logger
	supervisor *runtime_supervisor.Supervisor
	server     *http.Server
}

func buildApp(log *zap.Logger, sys *SystemContext) (*App, error) {

	if sys.Exec == nil {
		return nil, errors.New("missing execution context")
	}

	// --- Runtime ---
	rtx, err := runtime_engine.Build(sys.Exec, sys.Session, log)
	if err != nil {
		return nil, err
	}

	// --- Modules ---
	registry := kernel_registry.DefaultRegistry()

	// filtering + capability enforcement should already be resolved in boot
	ordered, err := kernel_supervisor.ResolveDependencies(registry)
	if err != nil {
		return nil, err
	}

	modules := modules_adapter.AdaptModules(ordered, rtx)

	if len(modules) == 0 {
		return nil, errors.New("no modules available after adaptation")
	}

	// --- Supervisor ---
	sup := runtime_supervisor.NewSupervisor(log, modules)

	return &App{
		log:        log,
		supervisor: sup,
	}, nil
}

// ============================================================
// LIFECYCLE
// ============================================================

func (app *App) Start(ctx context.Context) error {

	if err := app.supervisor.Init(ctx); err != nil {
		return err
	}

	if err := app.supervisor.Start(ctx); err != nil {
		return err
	}

	app.startHTTP()
	return nil
}

func (app *App) Stop(ctx context.Context) error {
	if app.server != nil {
		_ = app.server.Shutdown(ctx)
	}
	return app.supervisor.Stop(ctx)
}
