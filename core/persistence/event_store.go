// core/persistence/event_store.go
package persistence

import (
	"context"
)

type EventStore interface {
	SaveEvent(context.Context, *ExternalEvent) error
}
