// internal/schema/environment/env_config.go
// This is the "Source of Truth" that everyone can safely import.

package internal_environment

import (
	"time"

	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
)

type EnvConfig struct {
	SchemaVersion int                        `json:"internal_version"`
	Discovery     DiscoveryProfile           `json:"discovery_profile"`
	GeneratedAt   time.Time                  `json:"generated_at"`
	Identity      MachineIdentity            `json:"identity"`
	Hardware      HardwareProfile            `json:"hardware"`
	Platform      PlatformResolution         `json:"platform"`
	Attestation   EnvAttestation             `json:"attestation"`
	EntityType    internal_common.EntityKind `json:"entity_type"`
	TierType      internal_common.TierType   `json:"tier_type"`
}
