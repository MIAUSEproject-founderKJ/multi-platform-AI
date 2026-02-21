// MIAUSEproject-founderKJ/multi-platform-AI/core/platform/boot.go

package platform

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

const (
	firstBootMarker      = "FirstBootMarker"
	lastKnownEnvKey      = "LastKnownEnv"
	currentSchemaVersion = 3
)

func RunBootSequence(v *security.IsolatedVault) (*schema.BootSequence, error) {

raw, err := probe.PassiveScan()
if err != nil {
	return nil, fmt.Errorf("passive scan failed: %w", err)
}

identity := &schema.MachineIdentity{
	MachineName: raw.InstanceID,
	Platform:    raw.PlatformType,
	OS:          raw.OS,
	Arch:        raw.Architecture,
}

if err := security.VerifyEnvironment(*identity); err != nil {
	return nil, err
}

	bm := &BootManager{
		Vault:    v,
		Identity: rawID,
	}

	// 3. First boot decision BEFORE any active scan
	return bm.DecideBootPath()
}