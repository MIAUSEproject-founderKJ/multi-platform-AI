//modules/inference_module.go performs AI inference and writes results to storage.
package modules

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/modules/inference"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
)

const (
	InferenceQueueSize = 5000
	InferenceWorkers   = 4
)


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

type InferenceModule struct {
	BaseModule

	model Model

	queue chan []byte

	workers sync.WaitGroup

	totalPredictions atomic.Uint64
	totalErrors      atomic.Uint64
}

func NewInferenceModule() DomainModule {

	m := &InferenceModule{
		BaseModule: BaseModule{
			name: "InferenceModule",
			deps: []string{"TelemetryModule"},
		},
		queue: make(chan []byte, InferenceQueueSize),
	}

	return m
}

func (m *InferenceModule) Init(ctx *runtime.ExecutionContext) error {

	m.InitBase(ctx)

	m.model = inference.NewHTTPModel("http://localhost:9000")

	m.logger.Info("InferenceModule initialized")

	return nil
}

func (m *InferenceModule) Run(ctx context.Context) error {

	m.setRunning(true)

	m.logger.Info("Inference workers starting")

	for i := 0; i < InferenceWorkers; i++ {

		m.workers.Add(1)

		go m.worker(ctx, i)
	}

	<-ctx.Done()

	close(m.queue)

	m.workers.Wait()

	m.setRunning(false)

	m.logger.Info("InferenceModule stopped")

	return nil
}

func (m *InferenceModule) worker(ctx context.Context, id int) {

	defer m.workers.Done()

	logger := m.logger.With(zap.Int("worker", id))

	for {

		select {

		case payload, ok := <-m.queue:

			if !ok {
				return
			}

			m.processEvent(ctx, payload, logger)

		case <-ctx.Done():
			return
		}
	}
}

func (m *InferenceModule) processEvent(
	ctx context.Context,
	payload []byte,
	logger *zap.Logger,
) {

	var event TelemetryEvent

	if err := json.Unmarshal(payload, &event); err != nil {
		m.totalErrors.Add(1)
		logger.Warn("invalid telemetry input", zap.Error(err))
		return
	}

	if m.model == nil {
		m.totalErrors.Add(1)
		logger.Warn("no inference model loaded")
		return
	}

	req := inference.PredictionRequest{
		DeviceID:  event.DeviceID,
		Timestamp: time.Unix(event.Timestamp, 0),
		Features: map[string]float64{
			"value": event.Value,
		},
	}

	result, err := m.model.Predict(ctx, req)

	if err != nil {
		m.totalErrors.Add(1)
		logger.Error("model prediction failed", zap.Error(err))
		return
	}

	infResult := InferenceResult{
		DeviceID:   result.DeviceID,
		Timestamp:  result.Timestamp.Unix(),
		Prediction: result.Confidence,
	}

	m.totalPredictions.Add(1)

	resultBytes, _ := json.Marshal(infResult)

	_ = m.ctx.Router.Publish("vehicle_control", resultBytes)
	_ = m.ctx.Router.Publish("database", resultBytes)
	_ = m.ctx.Router.Publish("audit", resultBytes)
}
result, err := m.model.Predict(ctx, req)

	if err != nil {

		m.totalErrors.Add(1)

		logger.Error("model prediction failed", zap.Error(err))

		return
	}

infResult := InferenceResult{
	DeviceID:   result.DeviceID,
	Timestamp:  result.Timestamp.Unix(),
	Prediction: result.Confidence,
}

	m.totalPredictions.Add(1)

	resultBytes, _ := json.Marshal(result)

	// publish results
	_ = m.ctx.Router.Publish("vehicle_control", resultBytes)
	_ = m.ctx.Router.Publish("database", resultBytes)
	_ = m.ctx.Router.Publish("audit", resultBytes)
}

func (m *InferenceModule) Handle(ctx context.Context, payload []byte) error {

	if len(payload) == 0 {
		return errors.New("empty inference payload")
	}

	select {

	case m.queue <- payload:
		return nil

	default:
		m.totalErrors.Add(1)
		return errors.New("inference queue full")
	}
}