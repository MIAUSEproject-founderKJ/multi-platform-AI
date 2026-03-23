// =========================
// PRODUCTION-GRADE MAIN.GO
// Refactored to align with:
// - Strict module lifecycle
// - RuntimeContext separation
// - Supervisor control
// - zap logging
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
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime"
)

// ============================================================
// ENTRYPOINT
// ============================================================

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	app, err := NewApp(logger)
	if err != nil {
		logger.Fatal("failed to initialize app", zap.Error(err))
	}

	if err := app.Start(ctx); err != nil {
		logger.Fatal("failed to start app", zap.Error(err))
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		logger.Error("shutdown incomplete", zap.Error(err))
	}
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
// CONSTRUCTOR
// ============================================================

func NewApp(logger *zap.Logger) (*App, error) {

	// --- BOOT PHASE (NO RUNTIME OBJECTS) ---
	vault, err := security.OpenVault()
	if err != nil {
		return nil, err
	}

	bootSeq, _, err := boot.RunBootSequence(vault)
	if err != nil {
		return nil, err
	}

	execCtx, err := boot.ResolveExecutionContext(bootSeq)
	if err != nil {
		return nil, err
	}

	// --- RUNTIME CONTEXT ---
	rtx, err := runtime.NewRuntimeContext(logger)
	if err != nil {
		return nil, err
	}

	// --- MODULE REGISTRATION ---
	registry := modules.DefaultRegistry()
	filtered := modules.FilterModules(registry, execCtx)
	ordered, err := modules.ResolveDependencies(filtered)
	if err != nil {
		return nil, err
	}

	// --- ADAPT MODULES TO NEW INTERFACE ---
	adapted := modules.AdaptModules(ordered, rtx)

	// --- SUPERVISOR ---
	sup := runtime.NewSupervisor(logger, adapted)

	// --- INIT MODULES ---
	if err := sup.Init(context.Background()); err != nil {
		return nil, err
	}

	return &App{
		log:        logger,
		supervisor: sup,
	}, nil
}

// ============================================================
// START
// ============================================================

func (a *App) Start(ctx context.Context) error {

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
// HTTP
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
