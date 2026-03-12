// modules/ingestion_module.go

package modules

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
)

const (
	MaxPayloadSize = 1024 * 1024 // 1MB
	QueueSize      = 5000
)

type TelemetryPayload struct {
	DeviceID  string  `json:"device_id"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type IngestionModule struct {
	ctx *boot.RuntimeContext

	queue chan TelemetryEvent

	workers sync.WaitGroup

	limiter *rate.Limiter

	healthy atomic.Bool

	// metrics counters
	totalRequests atomic.Uint64
	totalErrors   atomic.Uint64
}

func NewIngestionModule() DomainModule {
	return &IngestionModule{
		queue:   make(chan []byte, QueueSize),
		limiter: rate.NewLimiter(2000, 4000), // 2000 req/sec burst 4000
	}
}

func (m *IngestionModule) Name() string {
	return "IngestionModule"
}

func (m *IngestionModule) DependsOn() []string {
	return nil
}

func (m *IngestionModule) Init(ctx *boot.RuntimeContext) error {

	m.ctx = ctx

	m.healthy.Store(true)

	ctx.Logger.Info("IngestionModule initialized")

	return nil
}

func (m *IngestionModule) Run(ctx context.Context) error {

	m.ctx.Logger.Info("IngestionModule starting workers")

workerCount := runtime.NumCPU()

for i := 0; i < workerCount; i++ {

	id := i

	m.workers.Add(1)

	m.Go(ctx, "inference-worker", func() {

		m.worker(ctx, id)

	})
}

	<-ctx.Done()

	m.ctx.Logger.Info("IngestionModule shutting down")

	close(m.queue)

	m.workers.Wait()

	return nil
}

func (m *IngestionModule) worker(ctx context.Context, id int) {

	defer m.workers.Done()

	logger := m.ctx.Logger.With(zap.Int("worker", id))

	for {

		select {

		case payload, ok := <-m.queue:

			if !ok {
				return
			}

			func() {

				defer func() {
					if r := recover(); r != nil {

						logger.Error("worker panic recovered", zap.Any("panic", r))

						m.totalErrors.Add(1)
					}
				}()

				m.processPayload(ctx, payload)

			}()

		case <-ctx.Done():
			return
		}
	}
}

func (m *IngestionModule) processPayload(ctx context.Context, payload []byte) {

	start := time.Now()

	var data TelemetryPayload

	if err := json.Unmarshal(payload, &data); err != nil {

		m.totalErrors.Add(1)

		m.ctx.Logger.Warn("invalid telemetry payload",
			zap.Error(err),
		)

		return
	}

	if data.DeviceID == "" {

		m.totalErrors.Add(1)

		m.ctx.Logger.Warn("missing device_id")

		return
	}

	err := m.ctx.Router.Publish("telemetry", payload)

	if err != nil {

		m.totalErrors.Add(1)

		m.ctx.Logger.Error("failed to publish telemetry",
			zap.Error(err),
		)

		return
	}

	latency := time.Since(start)

	m.ctx.Logger.Debug("telemetry processed",
		zap.String("device", data.DeviceID),
		zap.Duration("latency", latency),
	)
}

func (m *IngestionModule) Handle(ctx context.Context, payload []byte) error {

	m.totalRequests.Add(1)

	if len(payload) == 0 {
		m.totalErrors.Add(1)
		return errors.New("empty payload")
	}

	if len(payload) > MaxPayloadSize {

		m.totalErrors.Add(1)

		return fmt.Errorf("payload exceeds max size %d", MaxPayloadSize)
	}

	if !m.limiter.Allow() {

		m.totalErrors.Add(1)

		return errors.New("rate limit exceeded")
	}

select {

case m.queue <- payload:

default:

	<-m.queue

	m.queue <- payload

	m.totalErrors.Add(1)

		return errors.New("ingestion queue full")
	}
}

func (m *IngestionModule) Healthy() bool {
	return m.healthy.Load()
}