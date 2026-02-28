//core/persistence/postgres_store.go

type PostgresStore struct {
    db *sql.DB
}

func (p *PostgresStore) SaveEvent(ctx context.Context, e *ExternalEvent) error {

    payload, _ := json.Marshal(e.CanonicalData)

    _, err := p.db.ExecContext(
        ctx,
        "INSERT INTO events (source, type, payload, confidence, timestamp) VALUES ($1,$2,$3,$4,$5)",
        e.Source,
        e.Type,
        payload,
        e.Confidence,
        e.Timestamp,
    )

    return err
}