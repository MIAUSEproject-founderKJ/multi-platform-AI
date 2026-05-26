// bootstrap/resolver/execution_context_resolver.go
package bootstrap_resolver

import (
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

type executionContextView struct {
	exec runtime_types.ExecutionContext
}

func (e *executionContextView) Platform() internal_environment.PlatformClass {
	return e.exec.PlatformClass
}

func (e *executionContextView) Capabilities() internal_environment.CapabilitySet {
	return e.exec.Capabilities
}

func (e *executionContextView) SecurityTier() user_setting.TrustLevel {
	return e.exec.TrustLevel
}

func (e *executionContextView) HasPermission(p user_setting.PermissionKey) bool {
	return e.exec.Permissions[p]
}

func (e *executionContextView) ServiceType() user_setting.ServiceType {
	return e.exec.Service
}

func ResolveExecutionContext(
	bootCtx runtime_types.ExecutionContext,
	session *user_setting.UserSession,
) runtime_types.ExecutionContext {

	return &executionContextView{
		exec: bootCtx,
	}
}
