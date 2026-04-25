// boot/runtime_types.go
package boot

import (
	"context"

	"go.uber.org/zap"

	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
)

// ExecutionContext is runtime-facing, derived from BootSequence.
// It MUST NOT contain raw BootSequence.
type ExecutionContext struct {
	Logger      *zap.Logger
	Session     *schema_identity.UserSession
	Permissions map[schema_identity.Permission]bool
	PermMask    schema_security.PermissionMask
	TrustLevel  schema_identity.TrustLevel
}

// Optional runtime component example
type AgentRuntime struct{}

func (a *AgentRuntime) Start(ctx context.Context) error {
	return nil
}

func (a *AgentRuntime) Stop(ctx context.Context) error {
	return nil
}
