// runtime/supervisor/runtime_supervisor.go
package supervisor

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

///////////////////////////////////////////////////////////////
// MODULE INTERFACE (EXPECTED CONTRACT)
///////////////////////////////////////////////////////////////

type Module interface {
	Name() string
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Health() error
}

type HealthStatus struct {
	Healthy  bool
	Degraded bool
	Failed   int
	Total    int
}

func (s *Supervisor) HealthStatus() HealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	failed := 0
	for _, st := range s.modules {
		if !st.healthy {
			failed++
		}
	}

	return HealthStatus{
		Healthy:  failed == 0,
		Degraded: failed > 0,
		Failed:   failed,
		Total:    total,
	}
}

func (s *Supervisor) ModuleCount() int {
	return len(s.modules)
}

func (s *Supervisor) RestartFailed(ctx context.Context) error {
	for _, st := range s.modules {
		if err := st.module.Health(); err != nil {
			_ = st.module.Stop(ctx)
			_ = st.module.Start(ctx)
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////
// RESTART POLICY
///////////////////////////////////////////////////////////////

type RestartPolicy struct {
	MaxRestarts int
	Window      time.Duration
	BackoffMin  time.Duration
	BackoffMax  time.Duration
}

var DefaultRestartPolicy = RestartPolicy{
	MaxRestarts: 5,
	Window:      60 * time.Second,
	BackoffMin:  1 * time.Second,
	BackoffMax:  30 * time.Second,
}

///////////////////////////////////////////////////////////////
// INTERNAL STATE
///////////////////////////////////////////////////////////////

type moduleState struct {
	module Module

	running bool
	healthy bool

	restarts []time.Time
}

// /////////////////////////////////////////////////////////////
// SUPERVISOR
// /////////////////////////////////////////////////////////////

type Supervisor struct {
	log     *zap.Logger
	modules map[string]*moduleState
	order   []string
	policy  RestartPolicy
	mu      sync.RWMutex
	wg      sync.WaitGroup
}

///////////////////////////////////////////////////////////////
// CONSTRUCTOR
///////////////////////////////////////////////////////////////

func NewSupervisor(log *zap.Logger, mods []Module) *Supervisor {
	states := make(map[string]*moduleState)
	order := make([]string, 0, len(mods))

	for _, m := range mods {
		name := m.Name()
		states[name] = &moduleState{
			module: m,
		}
		order = append(order, name)
	}

	return &Supervisor{
		log:     log,
		modules: states,
		order:   order,
		policy:  DefaultRestartPolicy,
	}
}

///////////////////////////////////////////////////////////////
// INIT PHASE
///////////////////////////////////////////////////////////////

func (s *Supervisor) Init(ctx context.Context) error {
	for _, name := range s.order {
		m := s.modules[name].module

		s.log.Info("initializing module", zap.String("module", name))

		if err := m.Init(ctx); err != nil {
			return err
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////
// START PHASE
///////////////////////////////////////////////////////////////

func (s *Supervisor) Start(ctx context.Context) error {
	for _, name := range s.order {
		st := s.modules[name]

		s.wg.Add(1)
		go s.run(ctx, name, st)
	}
	return nil
}

func (s *Supervisor) run(ctx context.Context, name string, st *moduleState) {
	defer s.wg.Done()

	backoff := s.policy.BackoffMin

	for {
		if ctx.Err() != nil {
			s.log.Info("stopping module due to context cancellation", zap.String("module", name))
			_ = st.module.Stop(context.Background())
			return
		}

		s.log.Info("starting module", zap.String("module", name))

		s.mu.Lock()
		st.running = true
		s.mu.Unlock()

		err := st.module.Start(ctx)

		s.mu.Lock()
		st.running = false
		st.healthy = err == nil
		st.restarts = append(st.restarts, time.Now())
		s.mu.Unlock()

		if err != nil {
			s.log.Error("module crashed",
				zap.String("module", name),
				zap.Error(err),
			)
		}

		if !s.allowRestart(st) {
			s.log.Error("restart limit exceeded, stopping module permanently",
				zap.String("module", name),
			)
			return
		}

		time.Sleep(backoff)

		backoff *= 2
		if backoff > s.policy.BackoffMax {
			backoff = s.policy.BackoffMax
		}
	}
}

///////////////////////////////////////////////////////////////
// RESTART CONTROL
///////////////////////////////////////////////////////////////

func (s *Supervisor) allowRestart(st *moduleState) bool {
	now := time.Now()
	cutoff := now.Add(-s.policy.Window)

	valid := 0
	for _, t := range st.restarts {
		if t.After(cutoff) {
			valid++
		}
	}

	return valid <= s.policy.MaxRestarts
}

///////////////////////////////////////////////////////////////
// STOP PHASE
///////////////////////////////////////////////////////////////

func (s *Supervisor) Stop(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

///////////////////////////////////////////////////////////////
// HEALTH CHECK
///////////////////////////////////////////////////////////////

func (s *Supervisor) AllHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for name, st := range s.modules {
		if err := st.module.Health(); err != nil {
			s.log.Warn("module unhealthy",
				zap.String("module", name),
				zap.Error(err),
			)
			return false
		}
	}
	return true
}
