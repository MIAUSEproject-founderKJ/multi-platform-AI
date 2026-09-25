//runtime/types/execution_context.go

package runtime_types

import (
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type ExecutionContext interface {
	Platform() internal_common.PlatformClass
	Capabilities() internal_environment.CapabilitySet
	SecurityTier() user_setting.TrustLevel

	// optional (very useful)
	HasPermission(user_setting.PermissionKey) bool
	ServiceType() user_setting.ServiceType
}
