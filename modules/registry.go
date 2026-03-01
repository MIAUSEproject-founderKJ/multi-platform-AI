// modules/registry.go
package modules

func DefaultRegistry() []DomainModule {

	return []DomainModule{
		NewIngestionModule(),
		NewInferenceModule(),
		NewDatabaseSinkModule(),
		NewTelemetryModule(),
		NewIndustrialProtocolModule(),
		NewVehicleControlModule(),
		NewAuditModule(),
	}
}