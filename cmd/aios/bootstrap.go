// =========================
// PRODUCTION-GRADE MAIN.GO - cmd/aios/bootstrap.go
// Boot → ExecutionContext → RuntimeContext → Modules → Supervisor → Interfaces
// =========================

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot"
	boot_orchestrator "github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/orchestrator"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/resolver"
	security_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	modules "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
	registry "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/registry"
	cli "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/adapter"
	engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
	supervisor "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/supervisor"
	"go.uber.org/zap"
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
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("panic recovered", zap.Any("error", r))
		}
	}()

	// --- PHASE 1: BOOT ---
	sysCtx, err := BuildSystemContext()
	if err != nil {
		logger.Fatal("boot failure", zap.Error(err))
	}

	// --- PHASE 2: RUNTIME ---
	app, err := BuildRuntime(logger, sysCtx)
	if err != nil {
		logger.Fatal("runtime build failure", zap.Error(err))
	}

	// --- START ---
	go app.watchdog(ctx)

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

	logger.Info("boot_complete",
		zap.String("session_id", sysCtx.Session.SessionID),
		zap.String("mode", string(sysCtx.Session.Mode)),
		zap.Any("tier", sysCtx.Session.Tier),
	)
}

type SystemContext struct {
	Boot    *schema_boot.BootContext
	Exec    *boot.ExecutionContext
	Session *schema_identity.UserSession
}

func (a *App) Stop(ctx context.Context) error {
	if a.server != nil {
		_ = a.server.Shutdown(ctx)
	}
	return a.supervisor.Stop(ctx)
}

// ============================================================
// PHASE 1: BOOT
// ============================================================

func BuildSystemContext() (*SystemContext, error) {
	var lastErr error

	for i := 0; i < 3; i++ {
		ctx, err := attemptBoot()
		if err == nil {
			return ctx, nil
		}
		lastErr = err
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return nil, fmt.Errorf("boot failed after retries: %w", lastErr)
}

// ============================================================
// APP STRUCT
// ============================================================

type App struct {
	log        *zap.Logger
	supervisor *supervisor.Supervisor
	server     *http.Server

	degraded bool
}

// ============================================================
// PHASE 2: RUNTIME BUILD
// ============================================================

func BuildRuntime(logger *zap.Logger, sys *SystemContext) (*App, error) {

	// --- RUNTIME CONTEXT ---
	rtx, err := engine.NewRuntimeContext(logger)
	if err != nil {
		return nil, err
	}

	startCLIInput(rtx.Context, logger)

	// 🔥 attach session + config
	if sys.Session == nil {
		return nil, errors.New("nil session from boot")
	}

	rtx.Session = sys.Session

	if sys.Session.Config != nil {
		rtx.Config = sys.Session.Config
	}

	// --- MODULE GRAPH ---
	reg := registry.DefaultRegistry()

	filtered := modules.FilterModules(reg, sys.Boot)

	if len(filtered) == 0 {
		logger.Warn("no modules available, falling back to CLI-only mode")

	}

	ordered, err := registry.ResolveDependencies(filtered)
	if err != nil {
		return nil, err
	}

	adapted := modules.AdaptModules(ordered, rtx)

	if len(adapted) == 0 {
		logger.Warn("fallback to CLI module")

		adapted = []supervisor.Module{
			NewRecoverableModule(cli.NewCLIModule(), logger),
		}
	}

	resilient := make([]supervisor.Module, 0, len(adapted))

	for _, m := range adapted {
		resilient = append(resilient, NewRecoverableModule(m, logger))
	}

	// --- SUPERVISOR ---
	sup := supervisor.NewSupervisor(logger, resilient)

	return &App{
		log:        logger,
		supervisor: sup,
	}, nil
}

// ============================================================
// START
// ============================================================
func (r *recoverableModule) Health() error {
	return nil // or delegate
}

func (r *recoverableModule) Stop(ctx context.Context) error {
	return nil
}
func (r *recoverableModule) Init(ctx context.Context) error {
	return r.inner.Init(ctx)
}

func (r *recoverableModule) Name() string {
	return r.inner.Name()
}

func (r *recoverableModule) Start(ctx context.Context) error {
	backoff := time.Second

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !r.allowExecution() {
			time.Sleep(2 * time.Second)
			continue
		}

		err := r.safeRun(ctx)

		if err == nil {
			r.state = stateClosed
			r.failCount = 0
			return nil
		}

		r.predictor.Record()
		r.failCount++
		r.lastFailure = time.Now()

		// pre-failure signal
		if r.predictor.Rate() > r.predictor.threshold*0.7 {
			r.log.Warn("pre-failure signal detected")
		}

		if r.predictor.IsUnstable() {
			r.log.Error("instability spike → cooldown")
			time.Sleep(10 * time.Second)
		}

		if r.failCount >= 3 {
			r.state = stateOpen
			r.log.Warn("circuit breaker OPEN")
		}

		time.Sleep(backoff)
		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func (a *App) runFallbackLoop(ctx context.Context) {
	a.log.Warn("running in degraded mode")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			status := a.supervisor.HealthStatus()

			if status.Healthy && !status.Degraded {
				a.log.Info("system recovered from degraded mode")
				a.degraded = false
				return
			}

			a.log.Warn("system still degraded",
				zap.Int("failed_modules", status.Failed),
				zap.Int("total_modules", status.Total),
			)

		case <-ctx.Done():
			return
		}
	}
}

func attemptBoot() (*SystemContext, error) {
	vault, err := security_persistence.OpenVault()
	if err != nil {
		return nil, err
	}

	bootCtx := schema_boot.BootContext{
		Vault: vault,
	}

	bootSeq, session, err := boot_orchestrator.RunBootSequence(bootCtx)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("nil session")
	}

	bootSeq.UserSession = session

	bootCtxResolved, err := resolver.ResolveBootContext(bootSeq)
	if err != nil {
		return nil, err
	}

	execCtx, err := resolver.ResolveExecutionContext(bootSeq)
	if err != nil {
		return nil, err
	}

	return &SystemContext{
		Boot:    bootCtxResolved,
		Exec:    execCtx,
		Session: session,
	}, nil
}

func (a *App) watchdog(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			status := a.supervisor.HealthStatus()

			switch {
			case status.Failed == 0:
				continue

			case status.Failed < status.Total/2:
				a.log.Warn("early degradation")
				_ = a.supervisor.RestartFailed(ctx)

			case status.Failed >= status.Total/2:
				a.log.Error("major degradation")
				a.degraded = true
				a.activateSafeMode()

			case status.Failed == status.Total:
				a.log.Error("system collapse")
				a.activateSafeMode()
			}

		case <-ctx.Done():
			return
		}
	}
}

func (r *recoverableModule) safeRun(ctx context.Context) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			r.log.Error("module panic recovered", zap.Any("panic", rec))
			err = fmt.Errorf("panic: %v", rec)
		}
	}()
	return r.inner.Start(ctx)
}

func startCLIInput(ctx context.Context, logger *zap.Logger) {
	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Print("> ")
				input, err := reader.ReadString('\n')
				if err != nil {
					logger.Warn("stdin read error", zap.Error(err))
					continue
				}

				input = strings.TrimSpace(input)

				if input == "exit" {
					return
				}

				logger.Info("user_input", zap.String("input", input))
			}
		}
	}()
}

func (f *FailurePredictor) RecordError() {
	now := time.Now()
	f.events = append(f.events, now)

	// cleanup old
	cutoff := now.Add(-f.window)
	filtered := f.events[:0]
	for _, t := range f.events {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	f.events = filtered
}

type recoverableModule struct {
	inner supervisor.Module
	log   *zap.Logger

	// resilience
	predictor   *FailurePredictor
	state       circuitState
	failCount   int
	lastFailure time.Time
}

type circuitState int

const (
	stateClosed circuitState = iota
	stateOpen
	stateHalfOpen
)

func (a *App) activateSafeMode() {
	a.log.Warn("SAFE MODE ENABLED")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Print("[SAFE] > ")
				input, _ := reader.ReadString('\n')

				switch strings.TrimSpace(input) {
				case "status":
					fmt.Println("degraded")
				case "exit":
					return
				default:
					fmt.Println("limited command")
				}
			}
		}
	}()
}

func notifyUser(msg string) {
	fmt.Println("[SYSTEM]", msg)
}

func (a *App) Start(ctx context.Context) error {
	if err := a.supervisor.Init(ctx); err != nil {
		return fmt.Errorf("init failed: %w", err)
	}

	if err := a.supervisor.Start(ctx); err != nil {
		a.log.Error("supervisor start failed", zap.Error(err))
		a.degraded = true
		a.activateSafeMode()
		return nil
	}

	a.startHTTPServer()
	return nil
}

type FailurePredictor struct {
	window    time.Duration
	threshold float64 // errors per second
	events    []time.Time
}

func (f *FailurePredictor) Record() {
	now := time.Now()
	f.events = append(f.events, now)

	cutoff := now.Add(-f.window)
	i := 0
	for _, t := range f.events {
		if t.After(cutoff) {
			f.events[i] = t
			i++
		}
	}
	f.events = f.events[:i]
}

func (f *FailurePredictor) Rate() float64 {
	return float64(len(f.events)) / f.window.Seconds()
}

func (f *FailurePredictor) IsUnstable() bool {
	return f.Rate() > f.threshold
}

func (r *recoverableModule) allowExecution() bool {
	switch r.state {
	case stateOpen:
		if time.Since(r.lastFailure) > 15*time.Second {
			r.state = stateHalfOpen
			return true
		}
		return false
	case stateHalfOpen:
		return true
	case stateClosed:
		return true
	default:
		return true
	}
}

// ============================================================
// HTTP (EDGE ADAPTER)
// ============================================================

func (a *App) startHTTPServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		status := a.supervisor.HealthStatus()

		resp := map[string]interface{}{
			"healthy":  status.Healthy,
			"degraded": status.Degraded,
			"failed":   status.Failed,
			"total":    status.Total,
		}

		code := http.StatusOK
		if !status.Healthy {
			code = http.StatusServiceUnavailable
		}

		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(resp)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	a.server = server

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("http server error", zap.Error(err))
		}
	}()
}

func NewRecoverableModule(m supervisor.Module, log *zap.Logger) supervisor.Module {
	return &recoverableModule{
		inner: m,
		log:   log,
		predictor: &FailurePredictor{
			threshold: 5,
			window:    30 * time.Second,
		},
		state: stateClosed,
	}
}

type SystemGuard struct {
	failures int
	last     time.Time
}

func (g *SystemGuard) Trip() bool {
	if time.Since(g.last) < 30*time.Second {
		g.failures++
	} else {
		g.failures = 1
	}
	g.last = time.Now()

	return g.failures >= 5
}

func FallbackMinimal() []supervisor.Module {
	return []supervisor.Module{
		cli.NewCLIModule(),
	}
}
