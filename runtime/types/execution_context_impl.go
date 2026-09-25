// runtime/types/execution_context_impl.go
package runtime_types

import (
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

// resolvedExecutionContext is the concrete, immutable implementation of
// ExecutionContext produced once at boot. It is never mutated after
// construction — see docs/entry.md "BootContext vs RuntimeContext".
type resolvedExecutionContext struct {
	platform     internal_common.PlatformClass
	capabilities internal_environment.CapabilitySet
	trustLevel   user_setting.TrustLevel
	service      user_setting.ServiceType
	permissions  map[user_setting.PermissionKey]bool
}

// NewExecutionContext is the only constructor for ExecutionContext.
// Values are copied in, and the permission map is copied defensively so
// callers can't mutate it after handoff.
func NewExecutionContext(
	platform internal_common.PlatformClass,
	caps internal_environment.CapabilitySet,
	trust user_setting.TrustLevel,
	service user_setting.ServiceType,
	perms map[user_setting.PermissionKey]bool,
) ExecutionContext {
	cp := make(map[user_setting.PermissionKey]bool, len(perms))
	for k, v := range perms {
		cp[k] = v
	}
	return &resolvedExecutionContext{
		platform:     platform,
		capabilities: caps,
		trustLevel:   trust,
		service:      service,
		permissions:  cp,
	}
}

func (e *resolvedExecutionContext) Platform() internal_common.PlatformClass { return e.platform }
func (e *resolvedExecutionContext) Capabilities() internal_environment.CapabilitySet {
	return e.capabilities
}
func (e *resolvedExecutionContext) SecurityTier() user_setting.TrustLevel { return e.trustLevel }
func (e *resolvedExecutionContext) ServiceType() user_setting.ServiceType { return e.service }
func (e *resolvedExecutionContext) HasPermission(p user_setting.PermissionKey) bool {
	return e.permissions[p]
}
