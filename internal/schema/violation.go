package schema

import (
	"fmt"
	"unicode"

	"manala/internal/validator"

	"github.com/xeipuuv/gojsonschema"
)

func ResultErrorViolation(resultError gojsonschema.ResultError) validator.Violation {
	// Violation
	violation := validator.NewViolation(
		lcFirst(resultError.Description()),
	)

	// Path
	violation.Path = FieldPath(resultError.Field())

	switch resultError.(type) {
	case *gojsonschema.InvalidTypeError:
		expected := resultError.Details()["expected"].(string)
		actual := resultError.Details()["given"].(string)

		violation.Type = validator.InvalidType
		violation.Message = fmt.Sprintf("invalid type, expected %s, actual %s", expected, actual)
		violation.StructuredMessage = "invalid type"
		violation.Arguments = append(violation.Arguments,
			"expected", expected,
			"actual", actual,
		)
	case *gojsonschema.RequiredError:
		property := resultError.Details()["property"].(string)

		violation.Type = validator.Required
		violation.Message = fmt.Sprintf("missing %s property", property)
		violation.StructuredMessage = "missing property"
		violation.Property = property
	case *gojsonschema.AdditionalPropertyNotAllowedError:
		property := resultError.Details()["property"].(string)

		violation.Type = validator.AdditionalPropertyNotAllowed
		violation.Message = fmt.Sprintf("additional property %s is not allowed", property)
		violation.StructuredMessage = "additional property is not allowed"
		violation.Path = violation.Path.Join(property)
	case *gojsonschema.StringLengthGTEError:
		minimum := resultError.Details()["min"].(int)

		violation.Type = validator.StringGte
		violation.Message = fmt.Sprintf("string length must be greater than or equal to %d", minimum)
		violation.StructuredMessage = "too short string length"
		violation.Arguments = append(violation.Arguments,
			"minimum", minimum,
		)
	case *gojsonschema.StringLengthLTEError:
		maximum := resultError.Details()["max"].(int)

		violation.Type = validator.StringLte
		violation.Message = fmt.Sprintf("string length must be less than or equal to %d", maximum)
		violation.StructuredMessage = "too long string length"
		violation.Arguments = append(violation.Arguments,
			"maximum", maximum,
		)
	}

	return violation
}

func lcFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))

		return u + str[len(u):]
	}

	return ""
}
