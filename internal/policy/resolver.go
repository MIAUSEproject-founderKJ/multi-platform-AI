//internal/policy/resolver.go

package policy

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

func ResolvePermissions(ctx *schema.BootContext) map[schema.Permission]bool {

	perms := make(map[schema.Permission]bool)

	for _, p := range BasePolicies {

		if ctx.Capabilities.Has(p.RequireCaps) {

			for _, g := range p.Grant {
				perms[g] = true
			}

			for _, d := range p.Deny {
				perms[d] = false
			}
		}
	}

	return perms
}
