//core/datapipeline/engine.go

type PipelineEngine struct {
    decoder     Decoder
    normalizer  Normalizer
    reducer     DataErrorReducer
    validator   Validator
    router      *EventRouter
    store       persistence.EventStore
    monitor     telemetry.PerformanceMonitor
}

func (p *PipelineEngine) Process(raw []byte) {

    start := time.Now()

    e, err := p.decoder.Decode(raw)
    if err != nil {
        p.monitor.RecordError(err)
        return
    }

    if err := p.normalizer.Normalize(e); err != nil {
        p.monitor.RecordError(err)
        return
    }

    p.reducer.Reduce(e)

    if err := p.validator.Validate(e); err != nil {
        p.monitor.RecordError(err)
        return
    }

    if err := p.router.Route(e); err != nil {
        p.monitor.RecordError(err)
        return
    }

    if err := p.store.SaveEvent(context.Background(), e); err != nil {
        p.monitor.RecordError(err)
        return
    }

    p.monitor.RecordExecution("pipeline", time.Since(start))
}