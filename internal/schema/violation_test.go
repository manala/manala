package schema_test

import (
	"testing"

	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/schema"

	"github.com/stretchr/testify/suite"
	"github.com/xeipuuv/gojsonschema"
)

type ViolationSuite struct{ suite.Suite }

func TestViolationSuite(t *testing.T) {
	suite.Run(t, new(ViolationSuite))
}

func (s *ViolationSuite) TestResultErrorViolation() {
	internalError := &gojsonschema.InternalError{}
	internalError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))

	invalidTypeError := &gojsonschema.InvalidTypeError{}
	invalidTypeError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))
	invalidTypeError.SetDetails(gojsonschema.ErrorDetails{
		"expected": "foo",
		"given":    "bar",
	})

	requiredError := &gojsonschema.RequiredError{}
	requiredError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))
	requiredError.SetDetails(gojsonschema.ErrorDetails{
		"property": "foo",
	})

	additionalPropertyNotAllowedError := &gojsonschema.AdditionalPropertyNotAllowedError{}
	additionalPropertyNotAllowedError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))
	additionalPropertyNotAllowedError.SetDetails(gojsonschema.ErrorDetails{
		"property": "foo",
	})

	tests := []struct {
		test        string
		resultError gojsonschema.ResultError
		expected    schema.Violation
	}{
		{
			test:        "InternalError",
			resultError: internalError,
			expected: schema.Violation{
				Message:           "",
				StructuredMessage: "",
				Arguments:         []any(nil),
				Path:              path.Path(""),
				Property:          "",
			},
		},
		{
			test:        "InvalidTypeError",
			resultError: invalidTypeError,
			expected: schema.Violation{
				Message:           "invalid type, expected foo, actual bar",
				StructuredMessage: "invalid type",
				Arguments: []any{
					"expected", "foo",
					"actual", "bar",
				},
				Path:     path.Path(""),
				Property: "",
			},
		},
		{
			test:        "RequiredError",
			resultError: requiredError,
			expected: schema.Violation{
				Message:           "missing foo property",
				StructuredMessage: "missing property",
				Arguments:         []any(nil),
				Path:              path.Path(""),
				Property:          "foo",
			},
		},
		{
			test:        "AdditionalPropertyNotAllowedError",
			resultError: additionalPropertyNotAllowedError,
			expected: schema.Violation{
				Message:           "additional property foo is not allowed",
				StructuredMessage: "additional property is not allowed",
				Arguments:         []any(nil),
				Path:              path.Path("foo"),
				Property:          "",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			s.Equal(test.expected, schema.ViolationFrom(test.resultError))
		})
	}
}
