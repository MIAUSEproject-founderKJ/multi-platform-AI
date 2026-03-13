// core/reduction/reducer.go
package reduction

type Validator interface {
	Validate(*ExternalEvent) error
}
