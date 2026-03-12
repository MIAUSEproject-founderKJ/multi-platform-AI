//modules/database_sink_module.go
/*• Persist telemetry
• Write inference results
• Ensure durability*/

package modules

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
)

const (
	DBQueueSize = 5000
	DBWorkers   = 2
)

type DatabaseSinkModule struct {
	BaseModule

	db *sql.DB

	queue chan TelemetryEvent

	workers sync.WaitGroup

	totalWrites atomic.Uint64
	totalErrors atomic.Uint64
}

func NewDatabaseSinkModule() DomainModule {

	m := &DatabaseSinkModule{
		BaseModule: BaseModule{
			name: "DatabaseSinkModule",
			deps: []string{"TelemetryModule"},
		},
		queue: make(chan []byte, DBQueueSize),
	}

	return m
}

func (m *DatabaseSinkModule) Init(ctx *boot.RuntimeContext) error {

	m.InitBase(ctx)

	m.db = ctx.DB

	if m.db == nil {
		return errors.New("database not configured")
	}

	m.logger.Info("database sink initialized")

	return nil
}

func (m *DatabaseSinkModule) Run(ctx context.Context) error {

	m.setRunning(true)

	m.logger.Info("database sink workers starting")

	for i := 0; i < DBWorkers; i++ {

		m.workers.Add(1)

		m.Go(ctx, "inference-worker", func() {
	m.worker(ctx, i)
})
	}

	<-ctx.Done()

	close(m.queue)

	m.workers.Wait()

	m.logger.Info("database sink stopped")

	m.setRunning(false)

	return nil
}

func (m *DatabaseSinkModule) worker(ctx context.Context, id int) {

	defer m.workers.Done()

	logger := m.logger.Named("inference_worker").
	With(zap.Int("worker", id))

	for {

		select {

		case payload, ok := <-m.queue:

			if !ok {
				return
			}

			_, err := m.db.ExecContext(ctx,
				"INSERT INTO telemetry(data) VALUES(?)",
				string(payload),
			)

			if err != nil {

				m.totalErrors.Add(1)

				logger.Error("database insert failed",
					zap.Error(err),
				)

				continue
			}

			m.totalWrites.Add(1)

		case <-ctx.Done():
			return
		}
	}
}

func (m *DatabaseSinkModule) Handle(ctx context.Context, payload []byte) error {

	if len(payload) == 0 {
		return errors.New("empty payload")
	}

select {

case m.queue <- payload:

default:

	<-m.queue

	m.queue <- payload

	m.totalErrors.Add(1)
		return errors.New("database queue full")
	}
}