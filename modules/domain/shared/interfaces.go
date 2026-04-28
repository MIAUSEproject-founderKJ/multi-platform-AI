// modules/domain/shared/interfaces.go
package shared

import "context"

// Legacy module contract (existing system)
type DomainModule interface {
	Name() string
	Init(ctx context.Context) error
	Run(ctx context.Context) error
}

// Optional runtime injection
type RuntimeAware interface {
	SetRuntime(ctx any) // or *engine.RuntimeContext (preferred)
}
