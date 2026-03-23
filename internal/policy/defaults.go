//internal/policy/defaults.go

package policy

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

var BasePolicies = []Policy{

	// -------------------------------
	// Desktop Personal AI
	// -------------------------------
	{
		Name: "desktop_personal",

		RequireCaps: schema.CapNetwork,

		Grant: []schema.Permission{
			schema.PermBasicRuntime,
			schema.PermConfigEdit,
		},
	},

	// -------------------------------
	// Industrial / Vehicle Safety
	// -------------------------------
	{
		Name: "safety_critical",

		RequireCaps: schema.CapSafetyCritical,

		Grant: []schema.Permission{
			schema.PermHardwareIO,
			schema.PermDiagnostics,
		},

		Deny: []schema.Permission{
			schema.PermConfigEdit, // no runtime config changes
		},
	},

	// -------------------------------
	// Secure Device
	// -------------------------------
	{
		Name: "secure_enclave",

		RequireCaps: schema.CapSecureEnclave,

		Grant: []schema.Permission{
			schema.PermAdmin,
		},
	},
}
