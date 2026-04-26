package schema

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/serrors"

	"github.com/xeipuuv/gojsonschema"
)

type Violation struct {
	Message           string
	StructuredMessage string
	Arguments         []any
	Path              path.Path
	Property          string
}

func (violation Violation) StructuredError() serrors.Error {
	// Message
	message := violation.Message
	if violation.StructuredMessage != "" {
		message = violation.StructuredMessage
	}

	err := serrors.New(message)

	// Arguments
	err = err.With(violation.Arguments...)

	// Path
	if violation.Path != "" {
		err = err.With("path", violation.Path.String())
	}

	// Property
	if violation.Property != "" {
		err = err.With("property", violation.Property)
	}

	return err
}

type Violations []Violation

func (violations Violations) Errors() []error {
	errs := make([]error, len(violations))
	for i := range violations {
		errs[i] = errors.New(violations[i].Message)
	}

	return errs
}

func (violations Violations) StructuredErrors() []error {
	errs := make([]error, len(violations))
	for i := range violations {
		errs[i] = violations[i].StructuredError()
	}

	return errs
}

func ViolationFrom(resultError gojsonschema.ResultError) Violation {
	violation := Violation{
		Message: lcFirst(resultError.Description()),
	}

	// Path
	violation.Path = FieldPath(resultError.Field())

	switch resultError.(type) {
	case *gojsonschema.InvalidTypeError:
		expected := resultError.Details()["expected"].(string)
		actual := resultError.Details()["given"].(string)

		violation.Message = fmt.Sprintf("invalid type, expected %s, actual %s", expected, actual)
		violation.StructuredMessage = "invalid type"
		violation.Arguments = append(violation.Arguments,
			"expected", expected,
			"actual", actual,
		)
	case *gojsonschema.RequiredError:
		property := resultError.Details()["property"].(string)

		violation.Message = fmt.Sprintf("missing %s property", property)
		violation.StructuredMessage = "missing property"
		violation.Property = property
	case *gojsonschema.AdditionalPropertyNotAllowedError:
		property := resultError.Details()["property"].(string)

		violation.Message = fmt.Sprintf("additional property %s is not allowed", property)
		violation.StructuredMessage = "additional property is not allowed"
		violation.Path = violation.Path.Join(property)
	case *gojsonschema.StringLengthGTEError:
		minimum := resultError.Details()["min"].(int)

		violation.Message = fmt.Sprintf("string length must be greater than or equal to %d", minimum)
		violation.StructuredMessage = "too short string length"
		violation.Arguments = append(violation.Arguments,
			"minimum", minimum,
		)
	case *gojsonschema.StringLengthLTEError:
		maximum := resultError.Details()["max"].(int)

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
