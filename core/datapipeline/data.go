// core/datapipeline/data.go
package datapipeline

import "time"
type ExternalEvent struct {
    Source        string
    Type          string
    Timestamp     time.Time
    RawPayload    []byte
    CanonicalData map[string]interface{}
    Confidence    float64
    Metadata      map[string]string
}