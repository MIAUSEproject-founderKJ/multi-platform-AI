//modules/inference/model_adapter.go

package inference

type ModelAdapter struct {
	engine TensorEngine
}

func NewModelAdapter(engine TensorEngine) *ModelAdapter {
	return &ModelAdapter{
		engine: engine,
	}
}
