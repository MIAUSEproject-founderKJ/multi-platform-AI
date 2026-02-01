//MIAUSEproject-founderKJ/multi-platform-AI/configs/defaults/env_config.go

package defaults

import (
    "multi-platform-AI/configs/platforms"
    "multi-platform-AI/core/platform"
    "multi-platform-AI/core/security"
    "time"
)

type EnvConfig struct {
    SchemaVersion int       `json:"schema_version"`
    GeneratedAt   time.Time `json:"generated_at"`

    Identity    platforms.MachineIdentity   `json:"identity"`
    Hardware    platforms.HardwareProfile   `json:"hardware"`
    Platform    platform.PlatformResolution `json:"platform"`
    Attestation security.EnvAttestation     `json:"attestation"`
}