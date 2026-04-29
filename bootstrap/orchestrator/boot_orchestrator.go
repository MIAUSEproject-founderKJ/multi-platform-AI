// bootstrap/orchestrator/boot_orchestrator.go
package bootstrap_orchestrator

import (
	boot_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/phases"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

// RunBootSequence performs full bootstrap → verification → session creation
func RunBootSequence(ctx bootstrap.BootContext) (*internal_environment.BootSequence, *user_setting.UserSession, error) {

	discovery, err := boot_phase.PhaseDiscovery()
	if err != nil {
		return nil, nil, err
	}

	identity, err := boot_phase.PhaseIdentity(discovery)
	if err != nil {
		return nil, nil, err
	}

	bootSeq, err := boot_phase.PhaseBootResolution(ctx.Vault, identity)
	if err != nil {
		return nil, nil, err
	}

	// Merge capabilities (keep if needed)
	_ = bootSeq.Capabilities // avoid unused error OR remove entirely

	capsProfile := boot_phase.PhaseCapability()

	preSession, err := boot_phase.PhaseInterface(ctx, capsProfile)
	if err != nil {
		return nil, nil, err
	}

	session, err := boot_phase.PhaseAttestation(ctx.Vault, identity, bootSeq, preSession)
	if err != nil {
		return nil, nil, err
	}
	boot_phase.PhaseModules() // Load modules after attestation

	bootSeq.Env.Attestation.SessionToken = user_setting.UserIdentity

	return bootSeq, session, nil
}
