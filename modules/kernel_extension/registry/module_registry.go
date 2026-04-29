// modules/kernel_extension/registry/module_registry.go
package kernel_registry

import kernel_lifecycle "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/lifecycle"

func DefaultRegistry() []kernel_lifecycle.DomainModule {
	return []kernel_lifecycle.DomainModule{
		kernel_lifecycle.NewIngestionModule(),
		kernel_lifecycle.NewTelemetryModule(),
		kernel_lifecycle.NewInferenceModule(),
		kernel_lifecycle.NewDatabaseSinkModule(),
		kernel_lifecycle.NewIndustrialProtocolModule(),
		kernel_lifecycle.NewAuditModule(),
	}
}

// Responsible for receiving raw external data streams
// such as sensors, files, microphones, or network inputs.
// Converts them into normalized envelopes
