//core/optimization/optimizer.go

package optimization

import "time"

type Optimizer interface {
	Name() string

	// Algorithm policy
	PrecisionMode() PrecisionMode
	InferenceMode() InferenceMode
	BatchSize() int

	// Error reducer
	ShouldRetry(error) bool
	MaxRetries() int
	Backoff(attempt int) time.Duration

	// Performance review
	RecordLatency(time.Duration)
	RecordError(error)
	Evaluate() OptimizationReport
}

type PrecisionMode string
type InferenceMode string

const (
	PrecisionFull      PrecisionMode = "full"
	PrecisionReduced   PrecisionMode = "reduced"
	PrecisionAggressive PrecisionMode = "aggressive_pruned"
)

const (
	InferenceDeterministic InferenceMode = "deterministic"
	InferenceAdaptive      InferenceMode = "adaptive"
	InferenceHighThroughput InferenceMode = "throughput"
)

type OptimizationReport struct {
	AvgLatency   time.Duration
	ErrorRate    float64
	Recommendation string
}

type defaultOptimizer struct {
	platform PlatformClass

	precision PrecisionMode
	inference InferenceMode
	batch     int

	retries int

	latencies []time.Duration
	errors    int
	total     int
}

func NewDefaultOptimizer(p PlatformClass) Optimizer {

	opt := &defaultOptimizer{
		platform: p,
		retries:  3,
		batch:    1,
	}

	switch p {

	case PlatformVehicle:
		opt.precision = PrecisionAggressive
		opt.inference = InferenceDeterministic
		opt.batch = 1
		opt.retries = 1

	case PlatformIndustrial:
		opt.precision = PrecisionReduced
		opt.inference = InferenceDeterministic
		opt.batch = 2
		opt.retries = 2

	case PlatformPC:
		opt.precision = PrecisionFull
		opt.inference = InferenceAdaptive
		opt.batch = 8
		opt.retries = 3

	case PlatformCloud:
		opt.precision = PrecisionFull
		opt.inference = InferenceHighThroughput
		opt.batch = 32
		opt.retries = 5

	default:
		opt.precision = PrecisionFull
		opt.inference = InferenceAdaptive
		opt.batch = 4
	}

	return opt
}


func (o *defaultOptimizer) RecordLatency(d time.Duration) {
	o.latencies = append(o.latencies, d)
	o.total++
}

func (o *defaultOptimizer) RecordError(err error) {
	if err != nil {
		o.errors++
	}
}

func (o *defaultOptimizer) Evaluate() OptimizationReport {

	var total time.Duration
	for _, l := range o.latencies {
		total += l
	}

	avg := time.Duration(0)
	if len(o.latencies) > 0 {
		avg = total / time.Duration(len(o.latencies))
	}

	errorRate := 0.0
	if o.total > 0 {
		errorRate = float64(o.errors) / float64(o.total)
	}

	rec := "stable"

	if errorRate > 0.1 {
		rec = "reduce batch size"
	}

	if avg > 200*time.Millisecond {
		rec = "enable aggressive pruning"
	}

	return OptimizationReport{
		AvgLatency: avg,
		ErrorRate:  errorRate,
		Recommendation: rec,
	}
}