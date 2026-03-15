// boot/runtime_types.go
package boot

import (
	"context"
	"log/slog"
)

type ExecCtx struct {
	Logger *slog.Logger
}

type AgentRuntime struct{}

func (a *AgentRuntime) Start(ctx context.Context) error {
	return nil
}

func (a *AgentRuntime) Stop(ctx context.Context) error {
	return nil
}
