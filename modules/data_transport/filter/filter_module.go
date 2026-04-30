// modules/data_transport/filter/filter_module.go

package transport_filter

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
)

func FilterModules(all []domain_shared.DomainModule, ctx *bootstrap.BootContext) []domain_shared.DomainModule {

	var out []domain_shared.DomainModule

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
