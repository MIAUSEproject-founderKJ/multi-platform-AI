//core/datapipeline/data.go

type ExternalEvent struct {
    Source        string
    Type          string
    Timestamp     time.Time
    RawPayload    []byte
    CanonicalData map[string]interface{}
    Confidence    float64
    Metadata      map[string]string
}