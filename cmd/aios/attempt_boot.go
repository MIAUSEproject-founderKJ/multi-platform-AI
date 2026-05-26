// cmd/aios/attempt_boot.go
package main

import (
	bootstrap "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	bootstrap_orchestrator "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/orchestrator"
	bootstrap_resolver "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/resolver"
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
)

func attemptBoot() (*SystemContext, error) {
	vault, err := verification_persistence.OpenVault() //load local secured key
	if err != nil {
		return nil, err
	}

	// Stage 1: Boot orchestration
	bootCtx := bootstrap.NewBootContext(vault)

	bootSeq, session, err := bootstrap_orchestrator.RunBootSequence(bootCtx)
	if err != nil {
		return nil, err
	}

	execCtx, err := bootstrap_resolver.ResolveExecutionContext(bootCtx, session)
	if err != nil {
		return nil, err
	}

	return &SystemContext{
		Boot:    bootCtx,
		Runtime: execCtx,
		Session: session,
	}, nil
}
