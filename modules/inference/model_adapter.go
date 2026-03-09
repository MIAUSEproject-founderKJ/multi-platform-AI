//modules/inference/model_adapter.go

package inference

import (
	"context"
)

type ModelAdapter struct {
	engine TensorEngine
}

func NewModelAdapter(engine TensorEngine) *ModelAdapter {
	return &ModelAdapter{
		engine: engine,
	}
}