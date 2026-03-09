//modules/inference/feature_builder.go

func buildFeatures(req PredictionRequest) []float32 {

	return []float32{
		float32(req.Features["value"]),
	}
}