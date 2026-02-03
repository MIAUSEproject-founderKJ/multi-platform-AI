//MIAUSEproject-founderKJ/multi-platform-AI/configs/defaults/env_config.go

package defaults

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/configs/platforms"
)

type EnvConfig struct {
	SchemaVersion int       `json:"schema_version"`
	GeneratedAt   time.Time `json:"generated_at"`

	Identity    platforms.MachineIdentity   `json:"identity"`
	Hardware    platforms.HardwareProfile   `json:"hardware"`
	Platform    platform.PlatformResolution `json:"platform"`
	Attestation security.EnvAttestation     `json:"attestation"`
}
