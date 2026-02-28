//core/reduction/clamp_reducer.go

func (r *ClampReducer) Reduce(e *ExternalEvent) {

    if temp, ok := e.CanonicalData["temperature_c"].(float64); ok {
        if temp < -50 {
            e.CanonicalData["temperature_c"] = -50
            e.Confidence -= 0.1
        }
        if temp > 200 {
            e.CanonicalData["temperature_c"] = 200
            e.Confidence -= 0.1
        }
    }
}