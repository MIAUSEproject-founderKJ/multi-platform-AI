// runtime/types/runtime_types.go
package runtime_types

import (
	"context"

	"go.uber.org/zap"

	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

// ExecutionContext is runtime-facing, derived from BootSequence.
// It MUST NOT contain raw BootSequence.
type ExecutionContext struct {
	Logger      *zap.Logger
	Session     *user_setting.UserSession
	Permissions map[user_setting.PermissionKey]bool
	PermMask    internal_verification.PermissionMask
	TrustLevel  user_setting.TrustLevel
}

// Optional runtime component example
type AgentRuntime struct{}

func (a *AgentRuntime) Start(ctx context.Context) error {
	return nil
}

func (a *AgentRuntime) Stop(ctx context.Context) error {
	return nil
}
