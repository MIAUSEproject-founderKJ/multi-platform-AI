//internal\policy\policy_model.go

package policy

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

type Policy struct {
	Name string

	// Required capabilities
	RequireCaps schema.CapabilitySet

	// Explicit permissions granted
	Grant []schema.Permission

	// Explicit restrictions
	Deny []schema.Permission
}
