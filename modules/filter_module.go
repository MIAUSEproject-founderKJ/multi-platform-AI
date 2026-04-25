// modules/filter_module.go
package modules

import schema_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"

func FilterModules(all []DomainModule, ctx *schema_boot.BootContext) []DomainModule {

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
