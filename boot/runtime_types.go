// boot/runtime_types.go
package boot

import (
	"context"

	"go.uber.org/zap"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ExecutionContext is runtime-facing, derived from BootSequence.
// It MUST NOT contain raw BootSequence.
type ExecutionContext struct {
	Logger      *zap.Logger
	Session     *schema.UserSession
	Permissions map[schema.Permission]bool
	PermMask    schema.PermissionMask
	TrustLevel  schema.TrustLevel
}

// Optional runtime component example
type AgentRuntime struct{}

func (a *AgentRuntime) Start(ctx context.Context) error {
	return nil
}

func (a *AgentRuntime) Stop(ctx context.Context) error {
	return nil
}