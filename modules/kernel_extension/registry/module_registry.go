// modules/registry/module_registry.go
package registry

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"

func DefaultRegistry() []modules.DomainModule {
	return []modules.DomainModule{
		modules.NewIngestionModule(),
		modules.NewTelemetryModule(),
		modules.NewInferenceModule(),
		modules.NewDatabaseSinkModule(),
		modules.NewIndustrialProtocolModule(),
		modules.NewAuditModule(),
	}
}

// Responsible for receiving raw external data streams
// such as sensors, files, microphones, or network inputs.
// Converts them into normalized envelopes
