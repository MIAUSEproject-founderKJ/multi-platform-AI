//core/agent/runtime.go

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

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/optimization"
)

type AgentRuntime struct {
	router *router.Router
}

func NewAgentRuntime(r *router.Router) *AgentRuntime {
	return &AgentRuntime{
		router: r,
	}
}

type RawInput struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Source  string          `json:"source"`
}

func (a *AgentRuntime) Process(ctx context.Context, opt optimization.Optimizer, raw []byte) error {

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