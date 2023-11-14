package option

import (
	"manala/app"
	"manala/internal/path"
	"manala/internal/validator"
)

func NewPathedValidator(option app.RecipeOption) *PathedValidator {
	return &PathedValidator{
		option: option,
	}
}

type PathedValidator struct {
	option app.RecipeOption
}

func (validator PathedValidator) Validate(value any) (validator.Violations, error) {
	// Path
	optionPath := validator.option.Path()

	value, err := path.NewAccessor(
		optionPath,
		value,
	).
		Get()
	if err != nil {
		return nil, err
	}

	violations, err := validator.option.Validate(value)
	if err != nil {
		return nil, err
	}

	// Set violations path
	for i := range violations {
		violations[i].Path = optionPath
	}

	return violations, nil
}
