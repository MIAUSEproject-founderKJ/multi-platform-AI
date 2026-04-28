// internal/schema/system/device_identity.go

package internal_environment

import (
	"time"

	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type IdentityProfile struct {
	MachineID    string
	MachineName  string
	OS           string
	Architecture string
}

type FormFactor string

const (
	FormDesktop  FormFactor = "desktop"
	FormLaptop   FormFactor = "laptop"
	FormTablet   FormFactor = "tablet"
	FormPhone    FormFactor = "phone"
	FormHandheld FormFactor = "handheld"
)

type MachineIdentity struct {
	MachineID    string          `json:"machine_id"`
	PlatformType PlatformClass   `json:"platform_type"`
	Hostname     string          `json:"hostname"`
	OS           string          `json:"os"`
	Arch         string          `json:"arch"`
	Hardware     HardwareProfile `json:"hardware"`

	EntityType EntityKind            `json:"entity_type"`
	TierType   user_setting.TierType `json:"tier_type"`

	PasswordHash string    `json:"password_hash,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
