//internal\schema\boot\boot_sequence.go

package internal_boot

import (
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type BootSequence struct {
	Env          *internal_environment.EnvConfig
	Mode         BootMode
	Attested     bool
	Capabilities internal_environment.CapabilitySet
	Service      user_setting.ServiceType
	Entity       internal_common.EntityKind
	Tier         internal_common.TierType
	UserSession  *user_setting.UserSession
}
