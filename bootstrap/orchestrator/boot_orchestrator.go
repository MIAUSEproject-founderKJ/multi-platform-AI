// bootstrap/orchestrator/bootstrap_orchestrator.go
package bootstrap_orchestrator

import (
	bootstrap "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	bootstrap_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/phases"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

// RunBootSequence performs full bootstrap → verification → session creation

func RunBootSequence(
	bootctx bootstrap.BootContext,
) (
	*internal_boot.BootSequence,
	*user_setting.UserSession,
	error,
) {

	discovery, err := bootstrap_phase.PhaseDiscovery()
	if err != nil {
		return nil, nil, err
	}

	identity, err := bootstrap_phase.PhaseIdentity(discovery)
	if err != nil {
		return nil, nil, err
	}

	bootSeq, err := bootstrap_phase.PhaseBootResolution(identity)
	if err != nil {
		return nil, nil, err
	}

	capsProfile := bootstrap_phase.PhaseCapability()

	preSession, err := bootstrap_phase.PhaseInterface(
	// auth manager,
	// auth UI,
	)
	if err != nil {
		return nil, nil, err
	}

	session, err := bootstrap_phase.PhaseAttestation(
		identity,
		bootSeq,
		preSession,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := bootstrap_phase.PhaseRuntime(session); err != nil {
		return nil, nil, err
	}

	bootstrap_phase.PhaseModules()

	return bootSeq, session, nil
}
