//modules/base_module.go

package modules

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
)

type BaseModule struct {
	name string
	deps []string

	ctx    *runtime.RuntimeContext
	logger *zap.Logger

	healthy atomic.Bool
	running atomic.Bool

	startTime time.Time

	wg sync.WaitGroup

	eventsProcessed atomic.Uint64
	errorsTotal     atomic.Uint64
}

func (b *BaseModule) Name() string {
	return b.name
}

func (b *BaseModule) DependsOn() []string {
	return b.deps
}

func (b *BaseModule) InitBase(ctx *runtime.RuntimeContext) {

	b.ctx = ctx

	b.logger = ctx.Logger.With(
		zap.String("module", b.name),
	)

	b.startTime = time.Now()

	b.healthy.Store(true)

	b.logger.Info("module initialized")
}

func (b *BaseModule) Healthy() bool {
	return b.healthy.Load()
}

func (b *BaseModule) setHealthy(v bool) {
	b.healthy.Store(v)
}

func (b *BaseModule) setRunning(v bool) {
	b.running.Store(v)
}

func (b *BaseModule) Running() bool {
	return b.running.Load()
}