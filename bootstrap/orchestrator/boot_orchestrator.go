// bootstrap/orchestrator/bootstrap_orchestrator.go
package bootstrap_orchestrator

import (
	bootstrap_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/phases"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

// RunBootSequence performs full bootstrap → verification → session creation
func RunBootSequence(bootctx runtime_types.ExecutionContext) (*internal_environment.BootSequence, *user_setting.UserSession, error) {

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

	// Merge capabilities (keep if needed)
	_ = bootSeq.Capabilities // avoid unused error OR remove entirely

	capsProfile := bootstrap_phase.PhaseCapability()

	preSession, err := bootstrap_phase.PhaseInterface(bootctx, capsProfile)
	if err != nil {
		return nil, nil, err
	}

	session, err := bootstrap_phase.PhaseAttestation(identity, bootSeq, preSession)
	if err != nil {
		return nil, nil, err
	}
	bootstrap_phase.PhaseModules() // Load modules after attestation

	bootSeq.Env.Attestation.SessionToken = user_setting.UserIdentity

	return bootSeq, session, nil
}
