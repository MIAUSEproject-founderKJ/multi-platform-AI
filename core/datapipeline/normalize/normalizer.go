// core/datapipeline/normalize/normalizer.go
package normalize

type Normalizer interface {
	Normalize(*ExternalEvent) error
}
