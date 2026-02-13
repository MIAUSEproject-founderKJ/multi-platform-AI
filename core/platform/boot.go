// MIAUSEproject-founderKJ/multi-platform-AI/core/platform/boot.go

package platform

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/hmi"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

// BootManager handles the lifecycle of the AIOS startup
type BootManager struct {
	Vault    *security.IsolatedVault
	Identity *Identity
	HMIPipe  chan hmi.Update
}



func RunBootSequence(v *security.IsolatedVault) (*schema.BootSequence, error) {
    // 1. Initialize the BootManager with the Vault
    bm := &BootManager{
        Vault:    v,
        Identity: &Identity{}, // Identity starts empty
        HMIPipe:  make(chan hmi.Update, 10),
    }

    // 2. STAGE 0: Passive Scan (Hardware "Passport")
    // We do this first so bm.Identity has enough info for security checks
    rawId, err := probe.PassiveScan()
    if err != nil {
        return nil, fmt.Errorf("initial identity probe failed: %w", err)
    }
    bm.Identity.RawPassport = rawId

    // 3. THE HANDOFF: ManageBoot()
    // bm now has the Vault and the RawPassport. It can now decide 
    // whether to trust cached data or hammer the hardware.
    sequence, err := bm.ManageBoot()
    if err != nil {
        return nil, fmt.Errorf("platform verification failed: %w", err)
    }

    // 4. Return the verified sequence (Trust Score, Mode, etc.)
    return sequence, nil
}