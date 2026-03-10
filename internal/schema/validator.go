package schema

import (
	"github.com/manala/manala/internal/validator"

	"github.com/xeipuuv/gojsonschema"
)

func NewValidator(schema Schema) *Validator {
	return &Validator{
		schema: schema,
	}
}

type Validator struct {
	schema Schema
}

func (v *Validator) Validate(value any) (validator.Violations, error) {
	// Validate
	result, err := gojsonschema.Validate(
		gojsonschema.NewRawLoader(map[string]any(v.schema)),
		gojsonschema.NewGoLoader(value),
	)
	if err != nil {
		return nil, err
	}

	// Is valid ?
	if result.Valid() {
		return nil, nil
	}

	// Violations
	var violations validator.Violations
	for _, resultError := range result.Errors() {
		violations = append(violations, ResultErrorViolation(resultError))
	}

	return violations, nil
}
