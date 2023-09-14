package schema

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"manala/internal/path"
	"manala/internal/validator"
	"regexp"
	"unicode"
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
		violations = append(violations, v.violation(resultError))
	}

	return violations, nil
}

func (v *Validator) violation(resultError gojsonschema.ResultError) validator.Violation {
	// Violation
	violation := validator.NewViolation(
		v.lcFirst(resultError.Description()),
	)

	// Path
	violation.Path = v.path(resultError.Field())

	switch resultError.(type) {
	case *gojsonschema.InvalidTypeError:
		expected := resultError.Details()["expected"].(string)
		actual := resultError.Details()["given"].(string)

		violation.Type = validator.INVALID_TYPE
		violation.Message = fmt.Sprintf("invalid type, expected %s, actual %s", expected, actual)
		violation.StructuredMessage = "invalid type"
		violation.Arguments = append(violation.Arguments,
			"expected", expected,
			"actual", actual,
		)
	case *gojsonschema.RequiredError:
		property := resultError.Details()["property"].(string)

		violation.Type = validator.REQUIRED
		violation.Message = fmt.Sprintf("missing %s property", property)
		violation.StructuredMessage = "missing property"
		violation.Property = property
	case *gojsonschema.AdditionalPropertyNotAllowedError:
		property := resultError.Details()["property"].(string)

		violation.Type = validator.ADDITIONAL_PROPERTY_NOT_ALLOWED
		violation.Message = fmt.Sprintf("additional property %s is not allowed", property)
		violation.StructuredMessage = "additional property is not allowed"
		violation.Path = violation.Path.Join(property)
	case *gojsonschema.StringLengthGTEError:
		minimum := resultError.Details()["min"].(int)

		violation.Type = validator.STRING_GTE
		violation.Message = fmt.Sprintf("string length must be greater than or equal to %d", minimum)
		violation.StructuredMessage = "too short string length"
		violation.Arguments = append(violation.Arguments,
			"minimum", minimum,
		)
	case *gojsonschema.StringLengthLTEError:
		maximum := resultError.Details()["max"].(int)

		violation.Type = validator.STRING_LTE
		violation.Message = fmt.Sprintf("string length must be less than or equal to %d", maximum)
		violation.StructuredMessage = "too long string length"
		violation.Arguments = append(violation.Arguments,
			"maximum", maximum,
		)
	}

	return violation
}

func (v *Validator) lcFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

var validatorFieldPathRegex = regexp.MustCompile(`\.(\d+)`)

func (v *Validator) path(field string) path.Path {
	if field == gojsonschema.STRING_CONTEXT_ROOT {
		field = ""
	}

	// Index: foo.0 -> foo[0]
	field = validatorFieldPathRegex.ReplaceAllString(field, "[${1}]")

	return path.Path(field)
}
