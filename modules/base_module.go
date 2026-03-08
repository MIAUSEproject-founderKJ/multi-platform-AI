//modules/base_module.go

package modules

import (
	"sync/atomic"

	"go.uber.org/zap"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
)

type BaseModule struct {
	name string
	deps []string

	ctx    *runtime.ExecutionContext
	logger *zap.Logger

	healthy atomic.Bool
	running atomic.Bool
}

func (b *BaseModule) Name() string {
	return b.name
}

func (b *BaseModule) DependsOn() []string {
	return b.deps
}

func (b *BaseModule) InitBase(ctx *runtime.ExecutionContext) {

	b.ctx = ctx

	b.logger = ctx.Logger.With(
		zap.String("module", b.name),
	)

	b.healthy.Store(true)
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