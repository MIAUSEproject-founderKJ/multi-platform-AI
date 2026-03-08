//modules/inference/model.go

package inference

import "context"

type Model interface {
	Predict(ctx context.Context, req PredictionRequest) (*PredictionResult, error)
}