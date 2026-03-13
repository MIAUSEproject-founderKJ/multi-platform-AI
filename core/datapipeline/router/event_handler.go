// core/datapipeline/router/event_handler.go
package router

type EventHandler interface {
	Handle(*ExternalEvent) error
}
