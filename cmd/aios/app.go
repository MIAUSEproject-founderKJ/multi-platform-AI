// cmd/aios/app.go
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

type App struct {
	log        *zap.Logger
	supervisor *runtime_supervisor.Supervisor
	server     *http.Server
}

func buildApp(log *zap.Logger, sys *SystemContext) (*App, error) {
	if log == nil {
		return nil, errors.New("logger is required")
	}

	if sys == nil {
		return nil, errors.New("system context is required")
	}

	if sys.Execution == nil {
		return nil, errors.New("missing execution context")
	}

	// --------------------------------------------------------
	// Runtime
	// --------------------------------------------------------

	rtContainer, err := runtime_engine.Build(
		sys.Execution,
		sys.Session,
		log,
		sys.DB,
		sys.Router,
	)
	if err != nil {
		return nil, err
	}

	rtx := rtContainer.Context()
	if rtx == nil {
		return nil, errors.New("failed to create runtime context")
	}

	registry := kernel_registry.DefaultRegistry()

	ordered, err := kernel_supervisor.ResolveDependencies(registry)
	if err != nil {
		return nil, err
	}

	modules := modules_adapter.AdaptModules(ordered, rtx)

	if len(modules) == 0 {
		return nil, errors.New("no modules available after adaptation")
	}

	// --------------------------------------------------------
	// Supervisor
	// --------------------------------------------------------

	sup := runtime_supervisor.NewSupervisor(log, modules)

	return &App{
		log:        log,
		supervisor: sup,
	}, nil
}

func (app *App) Start(ctx context.Context) error {
	if app == nil {
		return errors.New("app is nil")
	}

	if app.supervisor == nil {
		return errors.New("runtime supervisor is nil")
	}

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
	if app == nil {
		return nil
	}

	if app.server != nil {
		_ = app.server.Shutdown(ctx)
	}

	if app.supervisor != nil {
		return app.supervisor.Stop(ctx)
	}

	return nil
}
