// MIAUSEproject-founderKJ/multi-platform-AI/core/platform/boot.go

package platform

import (
	"errors"
	"fmt"
	"multi-platform-AI/core/platform/classify"
	"multi-platform-AI/core/platform/degrade"
	"multi-platform-AI/core/platform/probe"
	"multi-platform-AI/core/security"
	"multi-platform-AI/core/security/attestation"
	"multi-platform-AI/internal/logging"
)

func RunBootSequence(v *security.IsolatedVault) (*BootSequence, error) {
	// 1. INITIAL CHECK: Is this the first time on this drive?
	// We check the marker before doing any expensive hardware IO.
	isFirstBoot := v.IsMissingMarker("FirstBootMarker")
	
	logging.Info("Phase 1: Initializing Boot Manager (FirstBoot: %v)", isFirstBoot)
id, _ := probe.PassiveScan()

mgr := &platform.BootManager{
    Vault:    v,
    Identity: id, // The Identity we just found
}

return mgr.ManageBoot()

