// modules/registry.go
package modules

func DefaultRegistry() []DomainModule {

	return []DomainModule{
		NewIngestionModule(),
		NewInferenceModule(),
		NewDatabaseSinkModule(),
		NewTelemetryModule(),
		NewIndustrialProtocolModule(),
		NewVehicleControlModule(), //control
		NewAuditModule(), //audit the module status
	}
}

		// Responsible for receiving raw external data streams
		// such as sensors, files, microphones, or network inputs.
		// Converts them into normalized envelopes