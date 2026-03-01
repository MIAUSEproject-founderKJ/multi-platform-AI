//modules/telemetry_module.go exports metrics to network.

type TelemetryModule struct {
	BaseModule
	client  TelemetryClient
	running atomic.Bool
}

//contructor
func NewTelemetryModule() *TelemtryModule{
	return &TelemetryModule{
		BaseModule: BaseModule{name: "telemetry"},
	}
}

//Dependencies
func (m *TelemetryModule) DependsOn() []string {
	return []string{"inference"}
}

//Platform restriction
func (m *TelemetryModule) SupportedPlatforms() []runtime.PlatformClass {
	return []runtime.PlatformClass{
		runtime.PlatformPC,
		runtime.PlatformCloud,
	}
}

//required capabilities
func (m *TelemetryModule) RequiredCapabilities() []string{
	return []string{"network"}
}

//If network capability is missing, FilterModules silently removes it.
func (m *TelemetryModule) Optional() bool {
	return true
}


//telemetry init
func (m *TelemetryModule) Init(ctx *runtime.ExecutionContext) error {
	m.ctx = ctx

	client, err := NewTelemetryClient()
	if err != nil {
		return err
	}

	m.client = client
	return nil
}

//telemetry start. This module exports the optimizer's performance review.
func (m *TelemetryModule) Start() error {
	m.running.Store(true)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for m.running.Load() {
			<-ticker.C

			report := m.ctx.Optimizer.Evaluate()

			err := m.client.Send(report)
			if err != nil {
				m.ctx.Optimizer.RecordError(err)
			}
		}
	}()

	return nil
}