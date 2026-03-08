//modules/inference/types.go

package inference

import "time"

type PredictionRequest struct {
	DeviceID  string             `json:"device_id"`
	Timestamp time.Time          `json:"timestamp"`
	Features  map[string]float64 `json:"features"`
}

type PredictionResult struct {
	DeviceID   string    `json:"device_id"`
	Prediction string    `json:"prediction"`
	Confidence float64   `json:"confidence"`
	Timestamp  time.Time `json:"timestamp"`
}