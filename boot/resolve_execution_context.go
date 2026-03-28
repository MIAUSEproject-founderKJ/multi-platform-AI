// boot/resolve_execution_context.go
package boot

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func ResolveExecutionContext(
	bs *schema.BootSequence,
) (*ExecutionContext, error) {

	// ------------------------------------------------------------
	// 1. Validation
	// ------------------------------------------------------------

	if bs == nil {
		return nil, fmt.Errorf("boot sequence is nil")
	}

	if bs.UserSession == nil {
		return nil, fmt.Errorf("missing user session")
	}

	if bs.Env == nil || !bs.Env.Attestation.Valid {
		return nil, fmt.Errorf("invalid attestation state")
	}

	// ------------------------------------------------------------
	// 2. Logger Initialization (runtime concern)
	// ------------------------------------------------------------

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	// ------------------------------------------------------------
	// 3. Derive Permissions (reuse BootContext logic)
	// ------------------------------------------------------------

	bootCtx, err := ResolveBootContext(bs)
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// 4. Construct ExecutionContext
	// ------------------------------------------------------------

	ctx := &ExecutionContext{
		Logger:      logger,
		Session:     bs.UserSession,
		Permissions: bootCtx.Permissions,
		TrustLevel:  bootCtx.TrustLevel,
	}

	return ctx, nil
}