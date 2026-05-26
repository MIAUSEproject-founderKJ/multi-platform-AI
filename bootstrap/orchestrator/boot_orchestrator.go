// bootstrap/orchestrator/bootstrap_orchestrator.go
package bootstrap_orchestrator

import (
	bootstrap "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	bootstrap_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/phases"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

// RunBootSequence performs full bootstrap → verification → session creation

func RunBootSequence(bootctx bootstrap.BootContext) (*internal_environment.BootSequence, *user_setting.UserSession, error) {
	//PhaseDiscovery scans bus, hardware, network, and environment to build a profile of the current system state. It gathers information about available resources, connected devices, network status, and other relevant environmental factors. This phase is crucial for understanding the context in which the system is operating and for making informed decisions in subsequent phases.
	discovery, err := bootstrap_phase.PhaseDiscovery()
	if err != nil {
		return nil, nil, err
	}
	//PhaseIdentity uses the information from the discovery phase to establish a unique identity for the machine. This may involve generating or retrieving a machine ID, determining the platform type, and collecting other relevant attributes that can be used to uniquely identify the machine in future interactions.
	identity, err := bootstrap_phase.PhaseIdentity(discovery)
	if err != nil {
		return nil, nil, err
	}
	//PhaseBootResolution determines the appropriate boot path based on the machine's identity and current state. It decides whether to perform a cold boot, warm boot, or resume from a previous state. This phase may also involve checking for first boot conditions and marking them accordingly in the verification vault.
	bootSeq, err := bootstrap_phase.PhaseBootResolution(identity)
	if err != nil {
		return nil, nil, err
	}

	// Merge capabilities (keep if needed)
	_ = bootSeq.Capabilities // avoid unused error OR remove entirely
	//PhaseCapability confirms the capabilities of the system based on the PhaseDiscovery and the PhaseIdentity return. It assesses the available resources, hardware features, and software capabilities to determine what functionalities can be supported. This phase is essential for tailoring the system's behavior to its actual capabilities and for ensuring that subsequent operations are compatible with the system's limitations.
	capsProfile := bootstrap_phase.PhaseCapability()

	//PhaseInterface prepares the system for attestation by setting up necessary interfaces and pre-session state. It may involve initializing communication channels, preparing data structures, or performing any necessary setup that is required before the attestation process can begin.
	preSession, err := bootstrap_phase.PhaseInterface(capsProfile)
	if err != nil {
		return nil, nil, err
	}

	//PhaseAttestation performs the attestation process to verify the user's identity to establish a secure session. It uses the machine's identity, the boot sequence information, and the pre-session state to authenticate the user and create a session token. This phase is critical for ensuring that only authorized users can access the system and for establishing a secure context for future interactions.
	session, err := bootstrap_phase.PhaseAttestation(identity, bootSeq, preSession)
	if err != nil {
		return nil, nil, err
	}

	//PhaseModules loads the necessary modules based on the attestation results and the established session. It ensures that the appropriate software components are initialized and ready for use, tailored to the authenticated user's permissions and the system's capabilities. This phase is essential for preparing the system for its intended operations while maintaining security and efficiency.
	bootstrap_phase.PhaseModules()

	bootSeq.Env.Attestation.SessionToken = user_setting.UserIdentity

	return bootSeq, session, nil
}
