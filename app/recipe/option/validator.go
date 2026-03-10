package option

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/validator"
)

type PathedValidator struct {
	option app.RecipeOption
}

func NewPathedValidator(option app.RecipeOption) *PathedValidator {
	return &PathedValidator{
		option: option,
	}
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
