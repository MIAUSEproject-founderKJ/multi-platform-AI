//runtime/session.go

/*Deterministic startup/shutdown
Clear ownership of resources
Cancellation propagation
Concurrency safety
Observability hooks
Backpressure + error isolation
No hidden goroutines*/

package runtime

import (
	"context"
	"errors"
	"sync"

	core_agent "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/agent"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

type SessionState int

const (
	SessionCreated SessionState = iota
	SessionStarting
	SessionRunning
	SessionStopping
	SessionStopped
)

type RuntimeSession struct {
	execCtx runtime_types.ExecutionContext
	agent   *core_agent.AgentRuntime

	stateMu sync.RWMutex
	state   SessionState

	rootCtx    context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	errCh      chan error
	shutdownCh chan struct{}

	startOnce sync.Once
	stopOnce  sync.Once
}

/*
Key Properties:
State machine with explicit transitions
Dedicated root context
WaitGroup to track goroutines
Error channel for async failures
Idempotent Start/Stop
*/
func NewRuntimeSession(bootctx runtime_types.ExecutionContext, agent *core_agent.AgentRuntime) *RuntimeSession {
	return &RuntimeSession{
		execCtx:    bootctx,
		agent:      agent,
		state:      SessionCreated,
		errCh:      make(chan error, 8),
		shutdownCh: make(chan struct{}),
	}
}

func (s *RuntimeSession) Start(parent context.Context) error {
	var startErr error

	s.startOnce.Do(func() {

		if !s.transition(SessionCreated, SessionStarting) {
			startErr = errors.New("invalid session state transition")
			return
		}

		s.rootCtx, s.cancel = context.WithCancel(parent)

		// Start agent runtime first
		if err := s.agent.Start(s.rootCtx); err != nil {
			startErr = err
			return
		}

		// Launch supervision loop
		s.wg.Add(1)
		go s.runSupervisor()

		s.transition(SessionStarting, SessionRunning)
	})

	return startErr
}

/*Agent starts before supervision loop
Goroutines are tracked
Safe idempotency*/

/*Supervisor Loop
Handles:
Async errors
Context cancellation
Panic recovery*/

func (s *RuntimeSession) runSupervisor() {
	defer s.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			s.errCh <- errors.New("session panic")
		}
	}()

	for {
		select {

		case <-s.rootCtx.Done():
			return

		case err := <-s.errCh:
			if err != nil {
				s.Stop(context.Background())
				return
			}
		}
	}
}

/*Cancels context first
Stops dependencies explicitly
Waits for all goroutines
Supports timeout via passed context*/

func (s *RuntimeSession) Stop(ctx context.Context) error {
	var stopErr error

	s.stopOnce.Do(func() {

		if !s.transition(SessionRunning, SessionStopping) {
			stopErr = errors.New("invalid session state for stop")
			return
		}

		// Cancel root context
		if s.cancel != nil {
			s.cancel()
		}

		// Stop agent runtime
		if err := s.agent.Stop(ctx); err != nil {
			stopErr = err
		}

		// Wait for all goroutines
		done := make(chan struct{})
		go func() {
			s.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-ctx.Done():
			stopErr = ctx.Err()
		}

		s.transition(SessionStopping, SessionStopped)
		close(s.shutdownCh)
	})

	return stopErr
}

// Ensures valid lifecycle flow.
func (s *RuntimeSession) transition(from, to SessionState) bool {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.state != from {
		return false
	}

	s.state = to
	return true
}
