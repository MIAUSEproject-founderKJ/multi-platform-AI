//core/datapipeline/router/event_handler.go
type EventHandler interface {
    Handle(*ExternalEvent) error
}