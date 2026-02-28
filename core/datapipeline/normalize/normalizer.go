//core/datapipeline/normalize/normalizer.go

type Normalizer interface {
    Normalize(*ExternalEvent) error
}