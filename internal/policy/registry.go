//internal\policy\registry.go

package policy

import (
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
)

var BasePolicies = []Policy{
	{
		Name:        "desktop_personal",
		RequireCaps: schema_security.CapNetwork,
		Grant: []schema_identity.Permission{
			schema_identity.PermBasicRuntime,
			schema_identity.PermConfigEdit,
		},
	},

	{
		Name:        "safety_critical",
		RequireCaps: schema_security.CapSafetyCritical,
		Grant: []schema_identity.Permission{
			schema_identity.PermHardwareIO,
			schema_identity.PermDiagnostics,
		},
		Deny: []schema_identity.Permission{
			schema_identity.PermConfigEdit,
		},
	},

	{
		Name:        "secure_enclave",
		RequireCaps: schema_security.CapSecureEnclave,
		Grant: []schema_identity.Permission{
			schema_identity.PermAdmin,
		},
	},
}
