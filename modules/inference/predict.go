//modules/inference/predict.go

package inference

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

httpClient: &http.Client{
    Timeout: 5 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:       100,
        IdleConnTimeout:    90 * time.Second,
        DisableCompression: false,
    },
}

type HTTPModel struct {
	endpoint   string
	httpClient *http.Client
	timeout    time.Duration
}

func NewHTTPModel(endpoint string) *HTTPModel {
	return &HTTPModel{
		endpoint: endpoint,
		timeout:  3 * time.Second,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (m *HTTPModel) Predict(
	ctx context.Context,
	payload PredictionRequest,
) (*PredictionResult, error) {

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("serialize payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		m.endpoint+"/predict",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prediction request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prediction service returned %d", resp.StatusCode)
	}

	var result PredictionResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode prediction: %w", err)
	}

	return &result, nil
}