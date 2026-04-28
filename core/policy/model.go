//internal/policy/model.go

package policy

import (
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

type Policy struct {
	Name string

	RequireCaps internal_verification.CapabilitySet

	Grant []user_setting.PermissionKey
	Deny  []user_setting.PermissionKey
}
