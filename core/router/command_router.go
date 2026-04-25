//core/router/command_handler_contract.go

/*The router is responsible for:
• Input validation
• Error reduction
• Message normalization
• Dispatching to domain modules*/

package router

import (
	"context"
	"errors"
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
	Dispatch(ctx context.Context, env Envelope) error
	Next(ctx context.Context) (interface{}, error)
}

type DefaultRouter struct {
	queue chan interface{}
}

func NewDefaultRouter() *DefaultRouter {
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

func (r *DefaultRouter) Dispatch(ctx context.Context, env Envelope) error {
	select {
	case r.queue <- env:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *DefaultRouter) Next(ctx context.Context) (interface{}, error) {
	select {
	case evt, ok := <-r.queue:
		if !ok {
			return nil, errors.New("router stopped")
		}
		return evt, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func New() Router {
	return NewDefaultRouter()
}
