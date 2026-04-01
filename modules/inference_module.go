// modules/inference_module.go performs AI inference and writes results to storage.
package modules

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime"
	"go.uber.org/zap"
)

// --- Constants ---
const (
	InferenceQueueSize = 5000
	InferenceWorkers   = 4
	BatchSize          = 16
)

// --- Telemetry and Inference Types ---
type TelemetryEvent struct {
	DeviceID  string  `json:"device_id"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type InferenceResult struct {
	DeviceID   string  `json:"device_id"`
	Timestamp  int64   `json:"timestamp"`
	Prediction float64 `json:"prediction"`
}

// --- Placeholder Model Interface ---
type Model interface {
	Predict(ctx context.Context, req PredictionRequest) (PredictionResult, error)
}

type PredictionRequest struct {
	DeviceID  string
	Timestamp time.Time
	Features  map[string]float64
}

type PredictionResult struct {
	DeviceID   string
	Timestamp  time.Time
	Confidence mathutil.Q16
	Prediction float64
}

// --- InferenceModule ---
type InferenceModule struct {
	BaseModule

	ctx     *schema.BootContext
	runtime *runtime.RuntimeContext // runtime reference
	logger  *zap.Logger
	running atomic.Bool

	model   Model
	queue   chan TelemetryEvent
	workers sync.WaitGroup

	totalPredictions atomic.Uint64
	totalErrors      atomic.Uint64
}

func NewInferenceModule() DomainModule {
	m := &InferenceModule{}
	m.SetName("InferenceModule")
	m.SetDeps([]string{"TelemetryModule"})
	return m
}

func (m *InferenceModule) RequiredCapabilities() schema.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}
func (m *InferenceModule) Optional() bool {
	return false
}

func (m *InferenceModule) SetRuntime(rtx *runtime.RuntimeContext) {
	m.runtime = rtx
}

// --- Base Init ---
func (m *InferenceModule) Init(ctx *schema.BootContext) error {
	if m.runtime == nil {
		return errors.New("runtime not set")
	}
	m.logger = m.runtime.Logger
	m.queue = make(chan TelemetryEvent, InferenceQueueSize)
	m.runtime.Bus.Subscribe("audio.features")
	return nil
}

// --- Run Loop ---
func (m *InferenceModule) Run(ctx context.Context) error {
	m.running.Store(true)
	m.logger.Info("Inference workers starting")

	for i := 0; i < InferenceWorkers; i++ {
		id := i
		m.workers.Add(1)
		go func() {
			defer m.workers.Done()
			m.worker(ctx, id)
		}()
	}

	<-ctx.Done()
	close(m.queue)
	m.workers.Wait()
	m.running.Store(false)
	m.logger.Info("InferenceModule stopped")
	return nil
}

func (m *InferenceModule) worker(ctx context.Context, id int) {
	for {
		select {
		case event, ok := <-m.queue:
			if !ok {
				return
			}

			batch := []PredictionRequest{m.convert(event)}
			for i := 1; i < BatchSize; i++ {
				select {
				case e := <-m.queue:
					batch = append(batch, m.convert(e))
				default:
					i = BatchSize
				}
			}
			m.processBatch(ctx, batch)
		case <-ctx.Done():
			return
		}
	}
}

func (m *InferenceModule) convert(event TelemetryEvent) PredictionRequest {
	return PredictionRequest{
		DeviceID:  event.DeviceID,
		Timestamp: time.Unix(event.Timestamp, 0),
		Features: map[string]float64{
			"value": event.Value,
		},
	}
}

func (m *InferenceModule) Start() error                               { m.running.Store(true); return nil }
func (m *InferenceModule) Stop() error                                { m.running.Store(false); return nil }
func (m *InferenceModule) SupportedPlatforms() []schema.PlatformClass { return nil }

func (m *InferenceModule) processBatch(ctx context.Context, batch []PredictionRequest) {
	for _, req := range batch {
		result, err := m.model.Predict(ctx, req)
		if err != nil {
			m.totalErrors.Add(1)
			m.logger.Error("model prediction failed", zap.Error(err))
			continue
		}

		infResult := InferenceResult{
			DeviceID:   result.DeviceID,
			Timestamp:  result.Timestamp.Unix(),
			Prediction: result.Prediction,
		}

		m.totalPredictions.Add(1)
		resultBytes, err := json.Marshal(infResult)
		if err != nil {
			m.totalErrors.Add(1)
			m.logger.Error("result serialization failed", zap.Error(err))
			continue
		}

		if m.runtime != nil {
			msg := runtime.Message{
				Topic: "vehicle_control",
				Data:  resultBytes,
			}
			m.runtime.Bus.Publish(msg)

			msg.Topic = "database"
			m.runtime.Bus.Publish(msg)

			msg.Topic = "audit"
			m.runtime.Bus.Publish(msg)
		}
	}
}

// --- DomainModule Methods ---
func (m *InferenceModule) Handle(ctx context.Context, payload []byte) error {
	if len(payload) == 0 {
		return errors.New("empty inference payload")
	}
	var event TelemetryEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	select {
	case m.queue <- event:
		return nil
	default:
		return errors.New("queue full")
	}
}

func (m *InferenceModule) Allowed(ctx *schema.BootContext) bool {
	return true
}

func (m *InferenceModule) Name() string {
	return m.name
}

func (m *InferenceModule) DependsOn() []string {
	return m.deps
}

func (m *InferenceModule) Category() ModuleCategory {
	return ModuleDomain
}

func (m *InferenceModule) IsRunning() bool {
	return m.running.Load()
}
