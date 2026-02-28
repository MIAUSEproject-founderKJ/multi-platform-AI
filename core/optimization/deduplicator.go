//core/optimization/deduplicator.go

type Deduplicator struct {
    cache map[string]time.Time
    ttl   time.Duration
}