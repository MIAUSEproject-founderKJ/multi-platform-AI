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
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	bootstrap_orchestrator "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/orchestrator"
	bootstrap_resolver "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/resolver"
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"

	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"

	transport_filter "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/data_transport/filter"
	kernel_registry "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/registry"
	kernel_supervisor "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/supervisor"

	interface_adapter "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/interface_adapter"
	engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
	runtime_supervisor "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/supervisor"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
	runtime_engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sysCtx, err := buildSystemContext()
	if err != nil {
		logger.Fatal("BOOT_FAILED", zap.Error(err))
	}

	app, err := buildRuntime(ctx, logger, sysCtx)
	if err != nil {
		logger.Fatal("RUNTIME_BUILD_FAILED", zap.Error(err))
	}

	if err := app.Start(ctx); err != nil {
		logger.Fatal("START_FAILED", zap.Error(err))
	}

	go app.watchdog(ctx)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_ = app.Stop(shutdownCtx)

	logger.Info("SYSTEM_EXIT",
		zap.String("user", sysCtx.Session.Identity.Username),
	)
}

// ============================================================
// SYSTEM CONTEXT
// ============================================================

type SystemContext struct {
	Boot    *bootstrap.BootContext
	Exec    *runtime_types.ExecutionContext
	Session *user_setting.UserSession
}

func buildSystemContext() (*SystemContext, error) {
	var last error

	for i := 0; i < 3; i++ {
		ctx, err := attemptBoot()
		if err == nil {
			return ctx, nil
		}
		last = err
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return nil, fmt.Errorf("boot failed: %w", last)
}

func attemptBoot() (*SystemContext, error) {
	vault, err := verification_persistence.OpenVault()
	if err != nil {
		return nil, err
	}

	bootCtx := bootstrap.BootContext{Vault: vault}

	bootSeq, session, err := bootstrap_orchestrator.RunBootSequence(bootCtx)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("nil session")
	}

	bootSeq.UserSession = session

	bootResolved, err := bootstrap_resolver.ResolveBootContext(bootSeq)
	if err != nil {
		return nil, err
	}

	execCtx, err := bootstrap_resolver.ResolveExecutionContext(bootSeq)
	if err != nil {
		return nil, err
	}

	return &SystemContext{
		Boot:    bootResolved,
		Exec:    execCtx,
		Session: session,
	}, nil
}

// ============================================================
// APP
// ============================================================

type App struct {
	log        *zap.Logger
	supervisor *runtime_supervisor.Supervisor
	server     *http.Server

	degraded bool
}

func buildRuntime(ctx context.Context, log *zap.Logger, sys *SystemContext) (*App, error) {

	if sys.Exec == nil {
		return nil, errors.New("missing execution context")
	}

	// 🔥 STRICT LINEAGE: Runtime derives from ExecutionContext
	rtx, err := runtime_engine.Build(sys.Exec, sys.Session, log)
if err != nil {
    return nil, err
}
	rtx.Session = sys.Session
	if sys.Session.Config != nil {
		rtx.Config = sys.Session.Config
	}

	startCLI(ctx, log)

	// --- MODULE GRAPH ---
	reg := kernel_registry.DefaultRegistry()

	filtered := transport_filter.FilterModules(reg, sys.Boot)

	// 🔥 Capability enforcement hook (critical for your architecture)
	filtered = enforceCapabilities(filtered, sys.Boot, log)

	ordered, err := kernel_supervisor.ResolveDependencies(filtered)
	if err != nil {
		return nil, err
	}

	adapted := modules_adapter.AdaptModules(ordered, rtx)

	if len(adapted) == 0 {
		log.Warn("NO_MODULES_FALLBACK")
		adapted = []runtime_supervisor.Module{
			interface_adapter.NewCLIModule(),
		}
	}

	resilient := wrapModules(adapted, log)

	sup := runtime_supervisor.NewSupervisor(log, resilient)

	return &App{
		log:        log,
		supervisor: sup,
	}, nil
}

// ============================================================
// START / STOP
// ============================================================

func (a *App) Start(ctx context.Context) error {

	if err := a.supervisor.Init(ctx); err != nil {
		return err
	}

	if err := a.supervisor.Start(ctx); err != nil {
		a.degraded = true
		go a.safeMode(ctx)
		return err
	}

	a.startHTTP()
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	if a.server != nil {
		_ = a.server.Shutdown(ctx)
	}
	return a.supervisor.Stop(ctx)
}

// ============================================================
// WATCHDOG
// ============================================================

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
				a.log.Warn("DEGRADED",
					zap.Int("failed", status.Failed),
					zap.Int("total", status.Total),
				)
				_ = a.supervisor.RestartFailed(ctx)

			default:
				a.log.Error("CRITICAL_DEGRADATION",
					zap.Int("failed", status.Failed),
					zap.Int("total", status.Total),
				)
				a.degraded = true
				go a.safeMode(ctx)
			}

		case <-ctx.Done():
			return
		}
	}
}

// ============================================================
// SAFE MODE
// ============================================================

func (a *App) safeMode(parent context.Context) {
	a.log.Warn("SAFE_MODE_ACTIVE")

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print("[SAFE]> ")
			in, _ := reader.ReadString('\n')

			switch strings.TrimSpace(in) {
			case "status":
				fmt.Println("degraded")
			case "exit":
				return
			default:
				fmt.Println("restricted")
			}
		}
	}
}

// ============================================================
// MODULE RESILIENCE
// ============================================================

func wrapModules(in []runtime_supervisor.Module, log *zap.Logger) []runtime_supervisor.Module {
	out := make([]runtime_supervisor.Module, 0, len(in))
	for _, m := range in {
		out = append(out, newRecoverable(m, log))
	}
	return out
}

type recoverable struct {
	inner runtime_supervisor.Module
	log   *zap.Logger

	failures int
	last     time.Time
}

func newRecoverable(m runtime_supervisor.Module, log *zap.Logger) runtime_supervisor.Module {
	return &recoverable{inner: m, log: log}
}

func (r *recoverable) Name() string                   { return r.inner.Name() }
func (r *recoverable) Init(ctx context.Context) error { return r.inner.Init(ctx) }
func (r *recoverable) Stop(ctx context.Context) error { return r.inner.Stop(ctx) }
func (r *recoverable) Health() error                  { return r.inner.Health() }

func (r *recoverable) Start(ctx context.Context) error {

	backoff := time.Second

	for {
		err := r.safeRun(ctx)
		if err == nil {
			r.failures = 0
			return nil
		}

		r.failures++
		r.last = time.Now()

		sleep := backoff + time.Duration(rand.Int63n(int64(backoff)))
		time.Sleep(sleep)

		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func (r *recoverable) safeRun(ctx context.Context) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			r.log.Error("MODULE_PANIC", zap.Any("panic", rec))
			err = fmt.Errorf("panic: %v", rec)
		}
	}()
	return r.inner.Start(ctx)
}

// ============================================================
// HTTP
// ============================================================

func (a *App) startHTTP() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		status := a.supervisor.HealthStatus()

		code := http.StatusOK
		if !status.Healthy {
			code = http.StatusServiceUnavailable
		}

		
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(status)
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	a.server = server

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("HTTP_ERROR", zap.Error(err))
		}
	}()
}

// ============================================================
// CLI
// ============================================================

func startCLI(ctx context.Context, log *zap.Logger) {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Print("> ")
				in, err := reader.ReadString('\n')
				if err != nil {
					continue
				}
				log.Info("INPUT", zap.String("cmd", strings.TrimSpace(in)))
			}
		}
	}()
}

// ============================================================
// ADAPTER HOOKS (PRODUCTION IMPLEMENTATION)
// ============================================================

// ModuleFactory defines a strict contract for module instantiation.
// Every module in registry SHOULD implement this.
type ModuleFactory interface {
	Name() string
	RequiredCapabilities() []string
	New(*runtime_engine.RuntimeContext) (runtime_supervisor.Module, error)
}

// Optional interface for modules that can self-declare compatibility
type CapabilityAware interface {
	IsCompatible(map[string]bool) bool
}



// ============================================================
// CAPABILITY ENFORCEMENT
// ============================================================

func enforceCapabilities(
	in []interface{},
	boot *bootstrap.BootContext,
	log *zap.Logger,
) []interface{} {

	// 🔥 Build capability matrix from BootContext
	capMatrix := buildCapabilityMatrix(boot)

	out := make([]interface{}, 0, len(in))

	for _, raw := range in {

		switch m := raw.(type) {

		// Preferred: explicit capability declaration
		case ModuleFactory:

			if !checkRequiredCapabilities(m.RequiredCapabilities(), capMatrix) {
				log.Warn("MODULE_CAPABILITY_REJECTED",
					zap.String("module", m.Name()),
					zap.Any("required", m.RequiredCapabilities()),
				)
				continue
			}

			out = append(out, m)

		// Secondary: self-checking modules
		case CapabilityAware:

			if !m.IsCompatible(capMatrix) {
				log.Warn("MODULE_SELF_REJECTED",
					zap.String("type", fmt.Sprintf("%T", raw)),
				)
				continue
			}

			out = append(out, raw)

		// Unknown modules → fail closed (IMPORTANT)
		default:
			log.Warn("MODULE_NO_CAPABILITY_DECLARATION_DROPPED",
				zap.String("type", fmt.Sprintf("%T", raw)),
			)
			continue
		}
	}

	return out
}

// ============================================================
// CAPABILITY MATRIX
// ============================================================

func buildCapabilityMatrix(boot *bootstrap.BootContext) map[string]bool {

	matrix := make(map[string]bool)

	// 🔥 You MUST map real BootContext signals here
	// Below is a safe baseline pattern

	// Example assumptions — adjust to your real BootContext schema:

	if boot == nil {
		return matrix
	}

	// --- Hardware / IO ---
	if boot.Platform != nil {
		matrix["platform_detected"] = true
		matrix["platform_class:"+boot.Platform.Class] = true
	}

	if boot.Devices != nil {

		for _, d := range boot.Devices {

			switch strings.ToLower(d.Type) {

			case "microphone":
				matrix["mic"] = true

			case "camera":
				matrix["camera"] = true

			case "lidar":
				matrix["lidar"] = true

			case "gpu":
				matrix["gpu"] = true

			case "network":
				matrix["network"] = true

			case "storage":
				matrix["storage"] = true
			}
		}
	}

	// --- Security / Session ---
	if boot.Vault != nil {
		matrix["secure_storage"] = true
	}

	// Extend:
	// matrix["robotics"]
	// matrix["vehicle"]
	// matrix["desktop"]
	// matrix["voice_enabled"]
	// matrix["vision_enabled"]

	return matrix
}

// ============================================================
// CAPABILITY CHECK
// ============================================================

func checkRequiredCapabilities(
	required []string,
	matrix map[string]bool,
) bool {

	if len(required) == 0 {
		return true
	}

	for _, cap := range required {
		if !matrix[cap] {
			return false
		}
	}
	return true
}