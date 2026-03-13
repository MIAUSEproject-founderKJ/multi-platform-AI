// core/validation/validator.go
package validation

type Validator interface {
	Validate(*ExternalEvent) error
}
