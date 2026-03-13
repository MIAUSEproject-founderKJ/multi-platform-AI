// core/decoder/json_decoder.go
package decoder

import (
	"encoding/json"
	"time"
)

type JSONDecoder struct{}

func (d *JSONDecoder) Decode(data []byte) (*ExternalEvent, error) {

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	return &ExternalEvent{
		RawPayload:    data,
		CanonicalData: payload,
		Timestamp:     time.Now(),
	}, nil
}
