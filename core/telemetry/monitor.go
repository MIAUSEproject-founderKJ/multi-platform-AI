//core/telemetry/monitor.go

type PerformanceMonitor struct {
    metrics map[string]Metric
    lock    sync.Mutex
}

func (p *PerformanceMonitor) RecordInference(duration time.Duration)
func (p *PerformanceMonitor) RecordExecution(domain string, duration time.Duration)
func (p *PerformanceMonitor) RecordError(err error)
func (p *PerformanceMonitor) Snapshot() Report

ctx.Monitor = telemetry.NewPerformanceMonitor()