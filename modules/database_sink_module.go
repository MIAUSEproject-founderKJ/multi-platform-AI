//modules/database_sink_module.go
/*• Persist telemetry
• Write inference results
• Ensure durability*/

package modules

import (
	"context"
	"database/sql"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
)

type DatabaseSinkModule struct {
	ctx     *runtime.ExecutionContext
	db      *sql.DB
	healthy atomic.Bool
}

func NewDatabaseSinkModule() DomainModule {
	return &DatabaseSinkModule{}
}

func (m *DatabaseSinkModule) Name() string {
	return "DatabaseSinkModule"
}

func (m *DatabaseSinkModule) DependsOn() []string {
	return []string{"TelemetryModule"}
}

func (m *DatabaseSinkModule) Init(ctx *runtime.ExecutionContext) error {

	m.ctx = ctx
	m.db = ctx.DB
	m.healthy.Store(true)

	return nil
}

func (m *DatabaseSinkModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *DatabaseSinkModule) Handle(ctx context.Context, payload []byte) error {

	_, err := m.db.ExecContext(ctx,
		"INSERT INTO telemetry(data) VALUES(?)",
		string(payload),
	)

	return err
}

func (m *DatabaseSinkModule) Healthy() bool {
	return m.healthy.Load()
}