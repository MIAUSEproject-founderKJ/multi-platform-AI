//internal\policy\model.go

package policy

import (
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_security "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/security"
)

type Policy struct {
	Name string

	RequireCaps schema_security.CapabilitySet

	Grant []schema_identity.Permission
	Deny  []schema_identity.Permission
}
