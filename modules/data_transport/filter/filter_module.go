// modules/filter_module.go
package modules

import internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"

func FilterModules(all []DomainModule, ctx *internal_boot.BootContext) []DomainModule {

	var out []DomainModule

	for _, m := range all {

		required := m.RequiredCapabilities()

		if !ctx.Capabilities.HasAll(required) {
			continue
		}

		if !m.Allowed(ctx) {
			continue
		}

		out = append(out, m)
	}

	return out
}
