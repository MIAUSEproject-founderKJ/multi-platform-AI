// modules/registry.go
package modules

func DefaultRegistry() []DomainModule {
	return []DomainModule{
		NewIngestionModule(),
		NewTelemetryModule(),
		NewInferenceModule(),
		NewDatabaseSinkModule(),
		NewIndustrialProtocolModule(),
		NewVehicleControlModule(),
		NewAuditModule(),
	}
}

		// Responsible for receiving raw external data streams
		// such as sensors, files, microphones, or network inputs.
		// Converts them into normalized envelopes