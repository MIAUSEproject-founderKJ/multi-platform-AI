//modules/inference_module.go performs AI inference and writes results to storage.
package modules

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"

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

	model inference.Model

	queue chan TelemetryEvent

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
		queue: make(chan TelemetryEvent, InferenceQueueSize),
	}

	return m
}

func (m *InferenceModule) Init(ctx *runtime.RuntimeContext) error {

	m.InitBase(ctx)

	engine := engines.NewONNXEngine()

err := engine.Load("models/model.onnx")
if err != nil {
	return err
}

adapter := inference.NewModelAdapter(engine)

m.model = adapter

	m.logger.Info("InferenceModule initialized")

	return nil
}

func (m *InferenceModule) Run(ctx context.Context) error {

	m.setRunning(true)

	m.logger.Info("Inference workers starting")

for i := 0; i < InferenceWorkers; i++ {

	id := i

	m.workers.Add(1)

	m.Go(ctx, "inference-worker", func() {
		m.worker(ctx, id)
	})
}

	<-ctx.Done()

	close(m.queue)

	m.workers.Wait()

	m.setRunning(false)

	m.logger.Info("InferenceModule stopped")

	return nil
}



/*waits for 1 event
tries to fill the batch without blocking
runs PredictBatch()*/
func (m *InferenceModule) worker(ctx context.Context, id int) {

	defer m.workers.Done()

	logger := m.logger.With(zap.Int("worker", id))

	for {

		select {

		case event, ok := <-m.queue:

			if !ok {
				return
			}

			batch := []inference.PredictionRequest{}

			req := m.convert(event)
			batch = append(batch, req)

			// collect additional events without blocking
			for i := 1; i < BatchSize; i++ {

				select {

				case e := <-m.queue:
					batch = append(batch, m.convert(e))

				default:
					break
				}
			}

			m.processBatch(ctx, batch, logger)

		case <-ctx.Done():
			return
		}
	}
}

func (m *InferenceModule) convert(event TelemetryEvent) inference.PredictionRequest {

	return inference.PredictionRequest{
		DeviceID:  event.DeviceID,
		Timestamp: time.Unix(event.Timestamp, 0),
		Features: map[string]float64{
			"value": event.Value,
		},
	}
}

func (b *BaseModule) LogInfo(msg string, fields ...zap.Field) {
	b.logger.Info(msg, fields...)
}

func (b *BaseModule) LogError(msg string, err error) {
	b.logger.Error(msg, zap.Error(err))
}

func (m *InferenceModule) processEvent(
	ctx context.Context,
	event TelemetryEvent,
	logger *zap.Logger,
) {

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

		logger.Error("model prediction failed",
			zap.Error(err),
		)

		return
	}

	infResult := InferenceResult{
		DeviceID:   result.DeviceID,
		Timestamp:  result.Timestamp.Unix(),
		Prediction: result.Confidence,
	}

	m.totalPredictions.Add(1)

	resultBytes, err := json.Marshal(infResult)

	if err != nil {

		m.totalErrors.Add(1)

		logger.Error("result serialization failed",
			zap.Error(err),
		)

		return
	}

	_ = m.ctx.Router.Publish("vehicle_control", resultBytes)
	_ = m.ctx.Router.Publish("database", resultBytes)
	_ = m.ctx.Router.Publish("audit", resultBytes)
}

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

		<-m.queue
		m.queue <- event

		m.totalErrors.Add(1)

		return errors.New("inference queue full")
	}
}

//Now workers are supervised automatically.
func (b *BaseModule) Go(ctx context.Context, name string, fn func()) {

	b.wg.Add(1)

	go func() {

		defer b.wg.Done()

		defer func() {
			if r := recover(); r != nil {

				b.errorsTotal.Add(1)

				b.logger.Error("worker panic recovered",
					zap.String("worker", name),
					zap.Any("panic", r),
				)
			}
		}()

		fn()
	}()
}

func (b *BaseModule) Shutdown() {

	b.logger.Info("shutting down module")

	b.wg.Wait()

	b.running.Store(false)

	b.logger.Info("module stopped")
}

func (b *BaseModule) IncEvents() {
	b.eventsProcessed.Add(1)
}

func (b *BaseModule) IncErrors() {
	b.errorsTotal.Add(1)
}

func (b *BaseModule) Stats() map[string]interface{} {

	return map[string]interface{}{
		"name":             b.name,
		"running":          b.running.Load(),
		"healthy":          b.healthy.Load(),
		"events_processed": b.eventsProcessed.Load(),
		"errors":           b.errorsTotal.Load(),
		"uptime":           time.Since(b.startTime).Seconds(),
	}
}

func (m *InferenceModule) processBatch(
	ctx context.Context,
	batch []inference.PredictionRequest,
	logger *zap.Logger,
) {

	results, err := m.model.PredictBatch(ctx, batch)

	if err != nil {

		m.totalErrors.Add(1)

		logger.Error("batch prediction failed", zap.Error(err))

		return
	}

	for _, result := range results {

		infResult := InferenceResult{
			DeviceID:   result.DeviceID,
			Timestamp:  result.Timestamp.Unix(),
			Prediction: result.Confidence,
		}

		resultBytes, err := json.Marshal(infResult)

		if err != nil {

			m.totalErrors.Add(1)

			logger.Error("result serialization failed",
				zap.Error(err),
			)

			continue
		}

		m.totalPredictions.Add(1)

		_ = m.ctx.Router.Publish("vehicle_control", resultBytes)
		_ = m.ctx.Router.Publish("database", resultBytes)
		_ = m.ctx.Router.Publish("audit", resultBytes)
	}
}

batchTimeout := time.After(5 * time.Millisecond)