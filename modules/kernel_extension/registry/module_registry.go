// modules/kernel_extension/registry/module_registry.go
package kernel_registry

import (
	transport_ingestion "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/data_transport/ingestion"
	transport_audit "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/data_transport/protocols"
	transport_storage "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/data_transport/storage"
	transport_telemetry "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/data_transport/telemetry"
	module_industrial "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/industrial"
	module_inference "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/inference"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
)

func DefaultRegistry() []domain_shared.DomainModule {
	return []domain_shared.DomainModule{
		transport_ingestion.NewIngestionModule(),
		transport_telemetry.NewTelemetryModule(),
		module_inference.NewInferenceModule(),
		transport_storage.NewDatabaseSinkModule(),
		module_industrial.NewIndustrialProtocolModule(),
		transport_audit.NewAuditModule(),
	}
}

// Responsible for receiving raw external data streams
// such as sensors, files, microphones, or network inputs.
// Converts them into normalized envelopes
