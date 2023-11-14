package validator

import "slices"

type Validator interface {
	Validate(value any) (Violations, error)
}

func New(opts ...Option) Validator {
	validator := &validator{}

	// Options
	for _, opt := range opts {
		opt(validator)
	}

	return validator
}

type validator struct {
	validators []Validator
	formatters []Formatter
}

func (validator *validator) Validate(value any) (Violations, error) {
	var violations Violations

	// Validators
	for _, validator := range validator.validators {
		_violations, err := validator.Validate(value)
		if err != nil {
			return nil, err
		}
		violations = append(violations, _violations...)
	}

	// Formatters
	for _, formatter := range validator.formatters {
		for i := range violations {
			formatter.Format(&violations[i])
		}
	}

	// Sort violations
	slices.SortStableFunc(violations, CompareViolations)

	return violations, nil
}

type Option func(validator *validator)

func WithValidators(validators ...Validator) Option {
	return func(validator *validator) {
		validator.validators = append(validator.validators, validators...)
	}
}

func Validators(vs ...Validator) Validator {
	return validators(vs)
}

type validators []Validator

func (validators validators) Validate(value any) (Violations, error) {
	var violations Violations
	for _, validator := range validators {
		_violations, err := validator.Validate(value)
		if err != nil {
			return nil, err
		}
		violations = append(violations, _violations...)
	}

	return violations, nil
}
