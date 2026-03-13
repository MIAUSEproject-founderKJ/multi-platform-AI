// core/datapipeline/router/event_router.go
package router

type EventRouter struct {
	handlers map[string][]EventHandler
}

func (r *EventRouter) Route(e *ExternalEvent) error {

	hs, ok := r.handlers[e.Type]
	if !ok {
		return nil
	}

	for _, h := range hs {
		if err := h.Handle(e); err != nil {
			return err
		}
	}

	return nil
}
