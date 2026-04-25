// boot/orchestrator/boot_orchestrator.go
package boot_orchestrator

import (
	boot_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/phases"
	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

// RunBootSequence performs full boot → verification → session creation
func RunBootSequence(ctx schema_boot.BootContext) (*schema_system.BootSequence, *schema_identity.UserSession, error) {

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

	bootSeq.Env.Attestation.SessionToken = session.SessionID

	return bootSeq, session, nil
}
