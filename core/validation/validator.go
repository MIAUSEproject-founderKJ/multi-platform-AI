//core/validation/validator.go
type Validator interface {
    Validate(*ExternalEvent) error
}