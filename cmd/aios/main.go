// cmd/aios/main.go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/optimization"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
	boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Start(ctx); err != nil {
		app.Logger.Fatal("app start failed", "error", err)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		app.Logger.Error("graceful shutdown incomplete", "error", err)
	}
}

type App struct {
	Ctx        *runtime.ExecutionContext
	Session    runtime.Session
	Modules    []modules.DomainModule
	Supervisor *Supervisor
	Watchdog   *Watchdog
	Logger     *slog.Logger
}

func NewApp() (*App, error) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	vault, err := security.OpenIsolatedVault()
	if err != nil {
		return nil, err
	}

	bootSeq, err := platform.RunBootSequence(vault)
	if err != nil {
		return nil, err
	}

	ctx, err := runtime.ResolveExecutionContext(bootSeq)
	if err != nil {
		return nil, err
	}

	ctx.Optimizer = optimization.NewDefaultOptimizer(ctx.PlatformClass)

	registry := modules.DefaultRegistry()
	filtered := modules.FilterModules(registry, ctx)

	ordered, err := modules.ResolveDependencies(filtered)
	if err != nil {
		return nil, err
	}

	app := &App{
		Ctx:     ctx,
		Modules: ordered,
		Logger:  logger,
	}

	app.Supervisor = NewSupervisor(logger)
	app.Watchdog = NewWatchdog(logger)

	return app, nil
}


func (a *App) Start(ctx context.Context) error {

	for _, m := range a.Modules {

		if err := m.Init(a.Ctx); err != nil {
			return err
		}

		a.Supervisor.Register(m)
	}

	if err := a.Supervisor.StartAll(ctx); err != nil {
		return err
	}

	router := NewDefaultRouter(a.Ctx)
	agent := NewAgentRuntime(router)

	a.Session = runtime.NewSession(a.Ctx, agent)

	if err := a.Session.Start(); err != nil {
		return err
	}

	a.Watchdog.Monitor(a.Modules)

	return nil
}


func (a *App) Stop(ctx context.Context) error {

	stopDone := make(chan struct{})

	go func() {
		a.Session.Stop()
		a.Supervisor.StopAll()
		close(stopDone)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopDone:
		return nil
	}
}


type Supervisor struct {
	modules map[string]modules.DomainModule
	logger  *slog.Logger
	retries int
}

func NewSupervisor(logger *slog.Logger) *Supervisor {
	return &Supervisor{
		modules: make(map[string]modules.DomainModule),
		logger:  logger,
		retries: 3,
	}
}

func (s *Supervisor) Register(m modules.DomainModule) {
	s.modules[m.Name()] = m
}

func (s *Supervisor) StartAll(ctx context.Context) error {
	for _, m := range s.modules {
		go s.runWithRecovery(ctx, m)
	}
	return nil
}

func (s *Supervisor) runWithRecovery(ctx context.Context, m modules.DomainModule) {
	attempt := 0

	for {
		err := safeStart(m)

		if err == nil {
			return
		}

		s.logger.Error("module crashed",
			"module", m.Name(),
			"error", err,
			"attempt", attempt,
		)

		attempt++

		if attempt > s.retries {
			s.logger.Error("module permanently disabled", "module", m.Name())
			return
		}

		time.Sleep(time.Second * 2)
	}
}

type HealthAware interface {
	Health() error
}

type Watchdog struct {
	logger *slog.Logger
}

func NewWatchdog(logger *slog.Logger) *Watchdog {
	return &Watchdog{logger: logger}
}

func (w *Watchdog) Monitor(mods []modules.DomainModule) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for _, m := range mods {
				if h, ok := m.(HealthAware); ok {
					if err := h.Health(); err != nil {
						w.logger.Warn("module unhealthy",
							"module", m.Name(),
							"error", err,
						)
					}
				}
			}
		}
	}()
}

func safeStart(m modules.DomainModule) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in %s: %v", m.Name(), r)
		}
	}()
	return m.Start()
}

logger := slog.New(handler).With("trace_id", traceID)

http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
	if !app.Supervisor.AllHealthy() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
})