// runtime/types/runtime_types.go
package runtime_types

import (
	"context"
)

type AgentRuntimeContext struct{}

func (a *AgentRuntimeContext) Start(ctx context.Context) error {
	return nil
}

func (a *AgentRuntimeContext) Stop(ctx context.Context) error {
	return nil
}
