// cmd/aios/attempt_boot.go
package main

import (
	"fmt"

	bootstrap "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	bootorchestrator "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/orchestrator"
	bootresolver "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/resolver"
	persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	"go.uber.org/zap"
)

func attemptBoot(log *zap.Logger) (*SystemContext, error) {
	vault, err := persistence.OpenVault(log)
	if err != nil {
		return nil, fmt.Errorf("boot failed at vault initialization: %w", err)
	}

	bootCtx := bootstrap.NewBootContext(vault)

	bootSeq, session, err := bootorchestrator.RunBootSequence(bootCtx)
	if err != nil {
		return nil, fmt.Errorf("boot sequence execution failed: %w", err)
	}

	execCtx, err := bootresolver.ResolveExecutionContext(bootSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve execution context: %w", err)
	}

	return &SystemContext{
		Boot:      bootCtx,
		Execution: execCtx,
		Session:   session,
	}, nil
}
