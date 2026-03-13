//core\router\default_router.go

/*The router is responsible for:
• Input validation
• Error reduction
• Message normalization
• Dispatching to domain modules*/

package router

import (
	"context"
)

type MessageType string

const (
	MessageTelemetry MessageType = "telemetry"
	MessageInference MessageType = "inference"
	MessageControl   MessageType = "control"
)

type Envelope struct {
	Type    MessageType
	Payload []byte
	Source  string
}

type Router interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Next(ctx context.Context) (interface{}, error)
}

type DefaultRouter struct {
	queue chan interface{}
}

func NewDefaultRouter(execCtx *ExecCtx) *DefaultRouter {
	return &DefaultRouter{
		queue: make(chan interface{}, 1024),
	}
}

func (r *DefaultRouter) Start(ctx context.Context) error {
	return nil
}

func (r *DefaultRouter) Stop(ctx context.Context) error {
	close(r.queue)
	return nil
}

func (r *DefaultRouter) Next(ctx context.Context) (interface{}, error) {
	select {
	case evt := <-r.queue:
		return evt, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
