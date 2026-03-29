// =========================
// PRODUCTION-GRADE MAIN.GO - cmd/aios/main.go
// Boot → ExecutionContext → RuntimeContext → Modules → Supervisor → Interfaces
// =========================

package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime"
)

// ============================================================
// ENTRYPOINT
// ============================================================

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, err := zap.NewProduction()
if err != nil {
	panic(err)
}	
	defer func() { _ = logger.Sync() }()

	// --- PHASE 1: BOOT ---
	sysCtx, err := BuildSystemContext()
	if err != nil {
		logger.Fatal("boot failure", zap.Error(err))
	}

	// --- PHASE 2: RUNTIME ---
	app, err := BuildRuntime(logger, sysCtx.Boot)
	if err != nil {
		logger.Fatal("runtime build failure", zap.Error(err))
	}

	// --- START ---
	if err := app.Start(ctx); err != nil {
		logger.Fatal("startup failure", zap.Error(err))
	}

	<-ctx.Done()

	// --- SHUTDOWN ---
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		logger.Error("shutdown incomplete", zap.Error(err))
	}
}

type SystemContext struct {
	Boot *schema.BootContext
	Exec *boot.ExecutionContext
	Session *schema.UserSession
}

// ============================================================
// PHASE 1: BOOT
// ============================================================

func BuildSystemContext() (*SystemContext, error) {

	vault, err := security.OpenVault()
	if err != nil {
		return nil, err
	}

	bootSeq, session, err := boot.RunBootSequence(vault)
	if err != nil {
		return nil, err
	}

	// Attach session back if not already embedded
	bootSeq.UserSession = session

	bootCtx, err := boot.ResolveBootContext(bootSeq)
	if err != nil {
		return nil, err
	}

	execCtx, err := boot.ResolveExecutionContext(bootSeq)
	if err != nil {
		return nil, err
	}

	return &SystemContext{
		Boot:    bootCtx,
		Exec:    execCtx,
		Session: session,
	}, nil
}

// ============================================================
// APP STRUCT
// ============================================================

type App struct {
	log        *zap.Logger
	supervisor *runtime.Supervisor
	server     *http.Server
}

// ============================================================
// PHASE 2: RUNTIME BUILD
// ============================================================

func BuildRuntime(logger *zap.Logger, sys *SystemContext) (*App, error) {

	// --- RUNTIME CONTEXT ---
	rtx, err := runtime.NewRuntimeContext(logger)
	if err != nil {
		return nil, err
	}

	// --- MODULE GRAPH ---
	registry := modules.DefaultRegistry()

	filtered := modules.FilterModules(registry, sys.Boot)

	ordered, err := modules.ResolveDependencies(filtered)
	if err != nil {
		return nil, err
	}

	adapted := modules.AdaptModules(ordered, rtx)

	// --- SUPERVISOR ---
	sup := runtime.NewSupervisor(logger, adapted)

	return &App{
		log:        logger,
		supervisor: sup,
	}, nil
}

// ============================================================
// START
// ============================================================

func (a *App) Start(ctx context.Context) error {

	// INIT is now explicitly part of runtime start phase
	if err := a.supervisor.Init(ctx); err != nil {
		return err
	}

	if err := a.supervisor.Start(ctx); err != nil {
		return err
	}

	a.startHTTPServer()

	return nil
}

// ============================================================
// STOP
// ============================================================

func (a *App) Stop(ctx context.Context) error {

	if a.server != nil {
		_ = a.server.Shutdown(ctx)
	}

	return a.supervisor.Stop(ctx)
}

// ============================================================
// HTTP (EDGE ADAPTER)
// ============================================================

func (a *App) startHTTPServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if !a.supervisor.AllHealthy() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	a.server = server

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("http server error", zap.Error(err))
		}
	}()
}

