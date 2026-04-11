package schema

import (
	"github.com/xeipuuv/gojsonschema"
)

type Validator struct {
	schema Schema
}

func NewValidator(schema Schema) *Validator {
	return &Validator{
		schema: schema,
	}
}

func (v *Validator) Validate(value any) (Violations, error) {
	// Validate
	result, err := gojsonschema.Validate(
		gojsonschema.NewRawLoader(map[string]any(v.schema)),
		gojsonschema.NewGoLoader(value),
	)
	if err != nil {
		return nil, err
	}

	// Is valid?
	if result.Valid() {
		return nil, nil
	}

	// Violations
	var violations Violations
	for _, resultError := range result.Errors() {
		violations = append(violations, ResultErrorViolation(resultError))
	}

	return violations, nil
}
