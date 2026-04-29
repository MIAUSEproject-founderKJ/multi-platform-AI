// internal/policy/registry.go
package policy

import (
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

var BasePolicies = []Policy{
	{
		Name:        "desktop_personal",
		RequireCaps: internal_environment.CapNetwork,
		Grant: []user_setting.PermissionKey{
			user_setting.PermBasicRuntime,
			user_setting.PermConfigEdit,
		},
	},

	{
		Name:        "safety_critical",
		RequireCaps: internal_environment.CapSafetyCritical,
		Grant: []user_setting.PermissionKey{
			user_setting.PermHardwareIO,
			user_setting.PermDiagnostics,
		},
		Deny: []user_setting.PermissionKey{
			user_setting.PermConfigEdit,
		},
	},

	{
		Name:        "secure_enclave",
		RequireCaps: internal_environment.CapSecureEnclave,
		Grant: []user_setting.PermissionKey{
			user_setting.PermAdmin,
		},
	},
}
