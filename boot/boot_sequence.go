//boot/boot_sequence.go

package boot

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// RunBootSequence performs full boot → verification → session creation
func RunBootSequence(ctx BootContext) (*schema.BootSequence, *schema.UserSession, error) {

	discovery, err := PhaseDiscovery()
	if err != nil {
		return nil, nil, err
	}

	identity, err := PhaseIdentity(discovery)
	if err != nil {
		return nil, nil, err
	}

	bootSeq, err := PhaseContext(ctx.Vault, identity)
	if err != nil {
    return nil, nil, err
	}

	// Merge capabilities
capSet := bootSeq.Capabilities

caps, err := DetectDeviceCapabilities(bootSeq.Env, capSet)
if err != nil {
    return nil, nil, err
}
bootSeq.Env.Discovery.Capabilities = *caps

preSession, err := PhaseAuthInterface(ctx, caps)
if err != nil { return nil, nil, err }

session, err := PhaseAttestation(ctx.Vault, identity, bootSeq)
if err != nil { return nil, nil, err }

	bootSeq.Env.Attestation.SessionToken = session.SessionID

	// 🔹 FULL INTERFACE
	err = PhaseMainInterface(session, caps)
	if err != nil {
		return nil, nil, err
	}

	return bootSeq, session, nil
}