// bootstrap/resolver/execution_context_resolver.go
package bootstrap_resolver

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

func ResolveExecutionContext(
	bs *internal_environment.BootSequence,
) (*runtime_types.ExecutionContext, error) {

	// ------------------------------------------------------------
	// 1. Validation
	// ------------------------------------------------------------

	if bs == nil {
		return nil, fmt.Errorf("bootstrap sequence is nil")
	}

	if bs.UserSession == nil {
		return nil, fmt.Errorf("missing user session")
	}

	if bs.Env == nil || !bs.Env.Attestation.Valid {
		return nil, fmt.Errorf("invalid attestation state")
	}

	// ------------------------------------------------------------
	// 2. Derive policy (pure logic)
	// ------------------------------------------------------------

	bootCtx, err := ResolveBootContext(bs) // must be PURE
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// 3. Construct ExecutionContext (NO runtime deps)
	// ------------------------------------------------------------

	ctx := &runtime_types.ExecutionContext{
		Session:     bs.UserSession,
		Permissions: bootCtx.Permissions,
		TrustLevel:  bootCtx.TrustLevel,
	}

	return ctx, nil
}
