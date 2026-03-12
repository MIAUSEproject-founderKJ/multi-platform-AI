//boot/boot_sequence.go

package boot

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/probe"
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


// Summary returns a human-readable log line for the console.
func (bs *BootSequence) Summary() string {
	return fmt.Sprintf("[%s] Booted as %s | Trust: %.0f%% | Mode: %s",
		bs.Timestamp.Format("15:04:05"),
		bs.PlatformID,
		bs.TrustScore*100,
		bs.Mode,
	)
}

// CanOperate returns true if the system allows any form of actuation.
func (bs *BootSequence) CanOperate() bool {
	return bs.Mode != "MANUAL_ONLY"
}

// IsAutonomous returns true only if the trust is high enough for self-governance.
func (bs *BootSequence) IsAutonomous() bool {
	return bs.Mode == "AUTONOMOUS"
}
