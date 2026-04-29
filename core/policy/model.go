//internal/policy/model.go

package policy

import (
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type Policy struct {
	Name string

	RequireCaps internal_environment.CapabilitySet

	Grant []user_setting.PermissionKey
	Deny  []user_setting.PermissionKey
}
