//pipeline_metrics.go
type PipelineMetrics struct {
    Throughput       float64
    AvgLatency       time.Duration
    ErrorRate        float64
    DBWriteLatency   time.Duration
}