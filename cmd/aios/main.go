// cmd/aios/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/agent"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/optimization"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
)

//////////////////////////////////////////////////////////////////
// ENTRYPOINT
//////////////////////////////////////////////////////////////////

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Start(rootCtx); err != nil {
		app.Logger.Error("app failed to start", "error", err)
		os.Exit(1)
	}

	<-rootCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		app.Logger.Error("shutdown incomplete", "error", err)
	}
}

//////////////////////////////////////////////////////////////////
// APP STRUCT
//////////////////////////////////////////////////////////////////

type App struct {
	Logger     *slog.Logger
	ExecCtx    *boot.RuntimeContext
	User       *schema.UserSession
	Session    *boot.Session
	Supervisor *Supervisor
	Server     *http.Server
}

//////////////////////////////////////////////////////////////////
// CONSTRUCTOR (BOOTSTRAP ONLY)
//////////////////////////////////////////////////////////////////

func NewApp() (*App, error) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	vault, err := security.OpenVault()
	if err != nil {
		return nil, err
	}

	bootSeq, userSession, err := boot.RunBootSequence(vault)
	if err != nil {
		return nil, err
	}

	execCtx, err := boot.ResolveExecutionContext(bootSeq)
	if err != nil {
		return nil, err
	}

	execCtx.Optimizer = optimization.NewDefaultOptimizer(execCtx.PlatformClass)

	registry := modules.DefaultRegistry()
	filtered := modules.FilterModules(registry, execCtx)
	ordered, err := modules.ResolveDependencies(filtered)
	if err != nil {
		return nil, err
	}

	supervisor := NewSupervisor(logger, ordered)

	app := &App{
		Logger:     logger,
		ExecCtx:    execCtx,
		User:       userSession,
		Supervisor: supervisor,
	}

	return app, nil
}

//////////////////////////////////////////////////////////////////
// START LIFECYCLE
//////////////////////////////////////////////////////////////////

func (a *App) Start(ctx context.Context) error {

	a.Logger.Info("initializing modules")

	//iterates over already dependency-sorted modules
	if err := a.Supervisor.InitAll(a.ExecCtx); err != nil {
		return err
	}

	a.Logger.Info("starting supervision tree")

	if err := a.Supervisor.Start(ctx); err != nil {
		return err
	}

	/*The router is responsible for:
	• Input validation
	• Error reduction
	• Message normalization
	• Dispatching to domain modules*/
	rtr := router.NewDefaultRouter(a.ExecCtx)

	/*• Algorithm distillation
	  • Optimization
	  • Confidence filtering
	  • Data shaping before dispatch
	*/
	agent := agent.NewAgentRuntime(rtr)

	/*The session handles:
	• External IO
	• Lifecycle binding
	• Controlled shutdown
	• Backpressure*/
	a.Session = boot.NewSession(a.ExecCtx, agent)

	if err := a.Session.Start(ctx); err != nil {
		return err
	}

	a.startHTTPServer()

	return nil
}

//////////////////////////////////////////////////////////////////
// STOP LIFECYCLE
//////////////////////////////////////////////////////////////////

func (a *App) Stop(ctx context.Context) error {

	a.Logger.Info("stopping application")

	if a.Server != nil {
		_ = a.Server.Shutdown(ctx)
	}

	if a.Session != nil {
		a.Session.Stop(ctx)
	}

	a.Supervisor.Stop()

	return nil
}

func (s *Supervisor) Stop() {
	s.wg.Wait()
}

//////////////////////////////////////////////////////////////////
// HTTP HEALTH ENDPOINTS
//////////////////////////////////////////////////////////////////

func (a *App) startHTTPServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if !a.Supervisor.AllHealthy() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	a.Server = server

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Error("http server error", "error", err)
		}
	}()
}

//////////////////////////////////////////////////////////////////
// SUPERVISOR
//////////////////////////////////////////////////////////////////

type Supervisor struct {
	logger  *slog.Logger
	modules []modules.DomainModule

	states map[string]*moduleState
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

type moduleState struct {
	module   modules.DomainModule
	healthy  bool
	started  bool
	restarts int
}

func NewSupervisor(logger *slog.Logger, mods []modules.DomainModule) *Supervisor {
	states := make(map[string]*moduleState)
	for _, m := range mods {
		states[m.Name()] = &moduleState{
			module:  m,
			healthy: false,
		}
	}
	return &Supervisor{
		logger:  logger,
		modules: mods,
		states:  states,
	}
}

func (s *Supervisor) InitAll(ctx *boot.RuntimeContext) error {
	for _, m := range s.modules {
		if err := m.Init(ctx); err != nil {
			return fmt.Errorf("module %s init failed: %w", m.Name(), err)
		}
	}
	return nil
}

func (s *Supervisor) Start(ctx context.Context) error {
	for _, m := range s.modules {
		s.wg.Add(1)
		go s.run(ctx, m)
	}
	return nil
}

func (s *Supervisor) run(ctx context.Context, m modules.DomainModule) {
	defer s.wg.Done()

	name := m.Name()
	backoff := time.Second

	for {

		func() {

			defer func() {
				if r := recover(); r != nil {
					s.logger.Error("module panic", "module", name, "panic", r)
				}
			}()

			s.mu.Lock()
			st := s.states[name]
			st.started = true
			st.healthy = true
			s.mu.Unlock()

			err := m.Run(ctx)

			if err != nil {
				s.logger.Error("module error", "module", name, "error", err)
			}

		}()

		if ctx.Err() != nil {
			return
		}

		s.mu.Lock()
		st := s.states[name]
		st.healthy = false
		st.restarts++
		s.mu.Unlock()

		time.Sleep(backoff)

		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func (s *Supervisor) AllHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, st := range s.states {
		if !st.started || !st.healthy {
			return false
		}
	}
	return true
}
