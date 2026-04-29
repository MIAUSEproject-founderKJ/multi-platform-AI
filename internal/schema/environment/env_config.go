// internal/schema/environment/env_config.go
// This is the "Source of Truth" that everyone can safely import.
package internal_environment

import (
	"time"

	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

type BootSequence struct {
	Env          *EnvConfig
	Mode         internal_boot.BootMode
	Attested     bool
	Capabilities internal_verification.CapabilitySet
	Service      user_setting.ServiceType
	Entity       EntityKind
	Tier         user_setting.TierType
	UserSession  *user_setting.UserSession
}

// Use EntityKind (uint8) internally for speed and clarity. Use TierType (string) externally for readability and compatibility.
type EntityKind uint8

const (
	EntityPersonal EntityKind = iota
	EntityOrganization
	EntityStranger
	EntityTester
)

func (m *MachineIdentity) BindHardware(env *EnvConfig) {
	m.Hardware = env.Hardware
}

type EnvConfig struct {
	SchemaVersion int                   `json:"internal_version"`
	Discovery     DiscoveryProfile      `json:"discovery_profile"`
	GeneratedAt   time.Time             `json:"generated_at"`
	Identity      MachineIdentity       `json:"identity"`
	Hardware      HardwareProfile       `json:"hardware"`
	Platform      PlatformResolution    `json:"platform"`
	Attestation   EnvAttestation        `json:"attestation"`
	EntityType    EntityKind            `json:"entity_type"`
	TierType      user_setting.TierType `json:"tier_type"`
}
