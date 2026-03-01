//core/router/router.go

/*The router is responsible for:
• Input validation
• Error reduction
• Message normalization
• Dispatching to domain modules*/

package router

import (
	"context"
	"errors"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
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

type Router struct {
	ctx      *runtime.ExecutionContext
	handlers map[MessageType]modules.DomainModule
}

func NewDefaultRouter(execCtx *runtime.ExecutionContext) *Router {
	r := &Router{
		ctx:      execCtx,
		handlers: make(map[MessageType]modules.DomainModule),
	}

	for _, m := range execCtx.ActiveModules {
		switch m.Name() {
		case "TelemetryModule":
			r.handlers[MessageTelemetry] = m
		case "InferenceModule":
			r.handlers[MessageInference] = m
		}
	}

	return r
}

func (r *Router) Dispatch(ctx context.Context, msg Envelope) error {

	if len(msg.Payload) == 0 {
		return errors.New("empty payload")
	}

	handler, ok := r.handlers[msg.Type]
	if !ok {
		return fmt.Errorf("no handler for message type %s", msg.Type)
	}

	return handler.Handle(ctx, msg.Payload)
}