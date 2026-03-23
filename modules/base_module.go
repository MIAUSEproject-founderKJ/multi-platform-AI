//modules/base_module.go

package modules

import (
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"go.uber.org/zap"
)

type BaseModule struct {
	name string
	deps []string

	// Add these fields for InferenceModule
	ctx     *schema.BootContext
	logger  *zap.Logger
	running atomic.Bool
}

func (b *BaseModule) setRunning(v bool) { b.running.Store(v) }
func (b *BaseModule) IsRunning() bool   { return b.running.Load() }

func (b *BaseModule) Init(ctx *schema.BootContext) {
	b.ctx = ctx
	b.logger = zap.NewExample() // replace with proper logger
}

func (b *BaseModule) LogInfo(msg string, fields ...zap.Field) {
	if b.logger != nil {
		b.logger.Info(msg, fields...)
	}
}

func (b *BaseModule) LogError(msg string, err error) {
	if b.logger != nil {
		b.logger.Error(msg, zap.Error(err))
	}
}

func (b *BaseModule) SetName(name string) {
	b.name = name
}

func (b *BaseModule) Name() string {
	return b.name
}

func (b *BaseModule) SetDeps(deps []string) {
	b.deps = deps
}

func (b *BaseModule) DependsOn() []string {
	return b.deps
}
