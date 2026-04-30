//modules/data_transport/storage/database_sink_module.go
/*• Persist telemetry
• Write inference results
• Ensure durability*/

package transport_storage

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
	kernel_lifecycle "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/kernel_extension/lifecycle"
	runtime_engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
	"go.uber.org/zap"
)

const (
	DBQueueSize = 5000
	DBWorkers   = 2
)

type DatabaseSinkModule struct {
	BaseModule
	runtime *runtime_engine.RuntimeContext
	db      *sql.DB
	queue   chan []byte

	workers sync.WaitGroup

	totalWrites atomic.Uint64
	totalErrors atomic.Uint64
	running     atomic.Bool

	logger *zap.Logger
}

// SetRuntime attaches the RuntimeContext
func (m *DatabaseSinkModule) SetRuntime(rtx *runtime_engine.RuntimeContext) {
	m.runtime = rtx
	m.db = rtx.DB
	m.logger = rtx.Logger
}

// Init subscribes to events
func (m *DatabaseSinkModule) Init(ctx *bootstrap.BootContext) error {
	if m.runtime == nil {
		return errors.New("runtime context not set")
	}

	m.runtime.Bus.Subscribe("database")
	return nil
}

// Run starts workers and handles shutdown
func (m *DatabaseSinkModule) Run(ctx context.Context) error {
	m.setRunning(true)

	for i := 0; i < DBWorkers; i++ {
		m.workers.Add(1)
		go m.worker(ctx, i)
	}

	<-ctx.Done()
	close(m.queue)
	m.workers.Wait()
	m.setRunning(false)
	return nil
}

// Handle enqueues payloads
func (m *DatabaseSinkModule) Handle(ctx context.Context, payload []byte) error {
	if len(payload) == 0 {
		return errors.New("empty payload")
	}

	select {
	case m.queue <- payload:
		return nil
	default:
		<-m.queue
		m.queue <- payload
		m.totalErrors.Add(1)
		return errors.New("database queue full")
	}
}

// worker routine
func (m *DatabaseSinkModule) worker(ctx context.Context, id int) {
	defer m.workers.Done()
	logger := m.logger.Named("db_worker").With(zap.Int("worker", id))

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
				logger.Error("database insert failed", zap.Error(err))
				continue
			}
			m.totalWrites.Add(1)
		case <-ctx.Done():
			return
		}
	}
}

// helper for running state
func (m *DatabaseSinkModule) setRunning(v bool) {
	m.running.Store(v)
}

// DomainModule implementation
func (m *DatabaseSinkModule) Name() string                            { return "DatabaseSinkModule" }
func (m *DatabaseSinkModule) Category() ModuleCategory                { return ModuleDomain }
func (m *DatabaseSinkModule) DependsOn() []string                     { return []string{"TelemetryModule"} }
func (m *DatabaseSinkModule) Allowed(ctx *bootstrap.BootContext) bool { return true }
func (m *DatabaseSinkModule) Start() error                            { return nil }
func (m *DatabaseSinkModule) Stop() error                             { return nil }
func (m *DatabaseSinkModule) SupportedPlatforms() []internal_environment.PlatformClass {
	return nil
}

// DomainModule implementation
func (m *DatabaseSinkModule) RequiredCapabilities() internal_environment.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}

func (m *DatabaseSinkModule) Optional() bool { return false }

// constructor
func NewDatabaseSinkModule() domain_shared.DomainModule {
	return &DatabaseSinkModule{
		BaseModule: kernel_lifecycle.BaseModule{name: "DatabaseSinkModule"},
		queue:      make(chan []byte, DBQueueSize),
	}
}
