package validator

type Validator interface {
	Validate(value any) (Violations, error)
}
