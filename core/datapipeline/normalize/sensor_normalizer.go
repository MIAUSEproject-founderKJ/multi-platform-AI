//core/datapipeline/normalize/sensor_normalizer.go

func (n *SensorNormalizer) Normalize(e *ExternalEvent) error {

    if v, ok := e.CanonicalData["temp"]; ok {
        e.CanonicalData["temperature_c"] = v
        delete(e.CanonicalData, "temp")
    }

    return nil
}