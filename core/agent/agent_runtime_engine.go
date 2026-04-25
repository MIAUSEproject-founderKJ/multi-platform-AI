//core/agent/agent_runtime_engine.go

/*The agent is responsible for:
• Algorithm distillation
• Optimization
• Confidence filtering
• Data shaping before dispatch*/

package agent

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
)

type AgentRuntime struct {
	router router.Router
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewAgentRuntime(router router.Router) *AgentRuntime {
	return &AgentRuntime{
		router: router,
	}
}

type RawInput struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Source  string          `json:"source"`
}

func (a *AgentRuntime) Process(ctx context.Context, opt Optimizer, raw []byte) error {

	var input RawInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return err
	}

	if input.Type == "" {
		return errors.New("missing message type")
	}

	// Optimization layer
	optimized, err := opt.Distill(input.Payload)
	if err != nil {
		return err
	}

	env := router.Envelope{
		Type:    router.MessageType(input.Type),
		Payload: optimized,
		Source:  input.Source,
	}

	return a.router.Dispatch(ctx, env)
}

func (a *AgentRuntime) Start(parent context.Context) error {
	a.ctx, a.cancel = context.WithCancel(parent)

	if err := a.router.Start(a.ctx); err != nil {
		return err
	}

	a.wg.Add(1)
	go a.eventLoop()

	return nil
}

func (a *AgentRuntime) Stop(ctx context.Context) error {
	if a.cancel != nil {
		a.cancel()
	}

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return ctx.Err()
	}

	return a.router.Stop(ctx)
}

func (a *AgentRuntime) eventLoop() {
	defer a.wg.Done()

	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			// Pull event from router
			event, err := a.router.Next(a.ctx)
			if err != nil {
				continue
			}

			// Process event
			a.handle(event)
		}
	}
}

func (a *AgentRuntime) handle(event interface{}) {
	// Domain-specific processing
}
