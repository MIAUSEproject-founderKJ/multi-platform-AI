// internal/policy/resolver.go
package policy

import (
	schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
)

type Resolver struct {
	Policies []Policy
}

func NewResolver() *Resolver {
	return &Resolver{
		Policies: BasePolicies,
	}
}

func (r *Resolver) Resolve(ctx *schema_boot.BootContext) Decision {

	perms := make(map[schema_identity.Permission]bool)

	var matched []string

	for _, p := range r.Policies {

		if ctx.Capabilities.Has(p.RequireCaps) {

			matched = append(matched, p.Name)

			for _, g := range p.Grant {
				perms[g] = true
			}

			for _, d := range p.Deny {
				perms[d] = false
			}
		}
	}

	decision := Decision{
		Granted: []string{},
		Denied:  []string{},
	}

	for perm, allowed := range perms {
		if allowed {
			decision.Granted = append(decision.Granted, string(perm))
		} else {
			decision.Denied = append(decision.Denied, string(perm))
		}
	}

	// simple rule: any denial overrides grant (safe default)
	for _, d := range decision.Denied {
		if d != "" {
			decision.Allowed = false
			decision.Reason = "explicit denial in policy: " + d
			return decision
		}
	}

	decision.Allowed = true
	decision.Reason = "policy matched: " + joinNames(matched)

	return decision
}

func joinNames(names []string) string {
	out := ""
	for i, n := range names {
		if i > 0 {
			out += ", "
		}
		out += n
	}
	return out
}
