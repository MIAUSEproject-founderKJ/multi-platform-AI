//core/reduction/reducer.go

type Validator interface {
    Validate(*ExternalEvent) error
}