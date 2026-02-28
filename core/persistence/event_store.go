//core/persistence/event_store.go

type EventStore interface {
    SaveEvent(context.Context, *ExternalEvent) error
}