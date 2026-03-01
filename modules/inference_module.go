//modules/inference_module.go performs AI inference and writes results to storage.
type InferenceModule struct {
	BaseModule
	model   ModelEngine
	running atomic.Bool
}

//constructor
func NewInferenceModule() *InferenceModule {
	return &InferenceModule{
		BaseModule: BaseModule{name: "inference"},
	}
}

//dependency declaration
func (m *InferenceModule) DependsOn() []string {
	return []string{"ingestion", "database_sink"}
}

//platform restriction, This prevents inference from running in cloud-only telemetry nodes.
func (m *InferenceModule) SupportedPlatforms() []runtime.PlatformClass {
	return []runtime.PlatformClass{
		runtime.PlatformVehicle,
		runtime.PlatformIndustrial,
		runtime.PlatformPC,
	}
}



func (m *InferenceModule) RequiredCapabilities() []string {
	return []string{"persistent_storage"}
}


//This ensures FilterModules panics if storage is missing.
func (m *InferenceModule) Optional() bool {
	return false
}

//Init Method (Where Optimizer Is Used)
/*This is algorithm distillation:
precision mode changes quantization
inference mode changes scheduling
batch size changes throughput vs latency*/
func (m *InferenceModule) Init(ctx *runtime.ExecutionContext) error {
	m.ctx = ctx

	precision := ctx.Optimizer.PrecisionMode()
	mode := ctx.Optimizer.InferenceMode()
	batch := ctx.Optimizer.BatchSize()

	engine, err := LoadModelEngine(precision, mode, batch)
	if err != nil {
		return err
	}

	m.model = engine
	return nil
}

/*This integrates:
performance review
error reduction
retry dampening
adaptive tuning hooks*/

func (m *InferenceModule) Start() error {
	m.running.Store(true)

	go func() {
		for m.running.Load() {

			start := time.Now()

			err := m.model.Run()

			m.ctx.Optimizer.RecordLatency(time.Since(start))
			m.ctx.Optimizer.RecordError(err)

			if err != nil && m.ctx.Optimizer.ShouldRetry(err) {
				time.Sleep(m.ctx.Optimizer.Backoff(1))
				continue
			}

			if err != nil {
				return
			}
		}
	}()

	return nil
}

func (m *InferenceModule) Stop() error {
	m.running.Store(false)
	return m.model.Close()
}