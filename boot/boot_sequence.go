//boot/boot_sequence.go

package boot

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// RunBootSequence performs full boot → verification → session creation
func RunBootSequence(v *security.IsolatedVault) (*schema.BootSequence, *schema.UserSession, error) {

	
	discovery, err := PhaseDiscovery()
	if err != nil {
		return nil, nil, err
	}

	identity, err := PhaseIdentity(discovery)
	if err != nil {
		return nil, nil, err
	}

	bootSeq, err := PhaseContext(v, identity)
	if err != nil {
		return nil, nil, err
	}

	session, err := PhaseAttestation(v, identity, bootSeq)
	if err != nil {
		return nil, nil, err
	}

	bootSeq.Env.Attestation.SessionToken = session.SessionID

	return bootSeq, session, nil
}
