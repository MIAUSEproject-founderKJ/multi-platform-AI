// modules/domain/shared/interfaces.go
package domain_shared

import runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"

// Legacy module contract (existing system)
type DomainModule interface {
	Name() string
	Init(ctx runtime_types.ExecutionContext) error
	Run(ctx runtime_types.ExecutionContext) error
}

// Optional runtime injection
type RuntimeAware interface {
	SetRuntime(ctx any) // or *runtime_engine.RuntimeContext (preferred)
}
