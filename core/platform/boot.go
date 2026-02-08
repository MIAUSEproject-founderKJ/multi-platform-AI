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

func RunBootSequence(v *security.IsolatedVault) (*BootManager, error) {
	isFirstBoot := v.IsMissingMarker("FirstBootMarker")
	logging.Info("Phase 1: Initializing Boot Manager (FirstBoot: %v)", isFirstBoot)

	// Perform initial probe
	rawId, _ := probe.PassiveScan()

	// Convert raw schema profile to our logic-capable Identity struct
	id := &Identity{IdentityProfile: rawId}

	mgr := &BootManager{
		Vault:    v,
		Identity: id,
		HMIPipe:  make(chan hmi.Update, 10),
	}

	return mgr, nil
}
