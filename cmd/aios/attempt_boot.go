// cmd/aios/attempt_boot.go
package main

import (
	"errors"

	bootstrap_orchestrator "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/orchestrator"
	bootstrap_resolver "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/resolver"
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

func attemptBoot() (*SystemContext, error) {

	vault, err := verification_persistence.OpenVault()
	if err != nil {
		return nil, err
	}

	// Stage 1: Boot orchestration
	bootSeq, session, err := bootstrap_orchestrator.RunBootSequence(
		runtime_types.ExecutionContext{Vault: vault},
	)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("nil session")
	}

	bootSeq.UserSession = session

	// Stage 2: Policy resolution (single source of truth)
	bootCtx, err := bootstrap_resolver.ResolveBootContext(bootSeq)
	if err != nil {
		return nil, err
	}

	// Stage 3: Execution projection
	execCtx, err := bootstrap_resolver.ResolveExecutionContext(bootCtx, session)
	if err != nil {
		return nil, err
	}

	return &SystemContext{
		Boot:    bootCtx,
		Exec:    execCtx,
		Session: session,
	}, nil
}
