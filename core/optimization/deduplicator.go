// core/optimization/deduplicator.go
package optimization

import (
	"time"
)

type Deduplicator struct {
	cache map[string]time.Time
	ttl   time.Duration
}
