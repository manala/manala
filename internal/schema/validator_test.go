package schema

import (
	"github.com/xeipuuv/gojsonschema"
	"manala/internal/path"
	"manala/internal/validator"
)

func (s *Suite) TestValidatorViolation() {
	internalError := &gojsonschema.InternalError{}
	internalError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))

	invalidTypeError := &gojsonschema.InvalidTypeError{}
	invalidTypeError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))
	invalidTypeError.SetDetails(gojsonschema.ErrorDetails{
		"expected": "expected",
		"given":    "given",
	})

	requiredError := &gojsonschema.RequiredError{}
	requiredError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))
	requiredError.SetDetails(gojsonschema.ErrorDetails{
		"property": "property",
	})

	additionalPropertyNotAllowedError := &gojsonschema.AdditionalPropertyNotAllowedError{}
	additionalPropertyNotAllowedError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))
	additionalPropertyNotAllowedError.SetDetails(gojsonschema.ErrorDetails{
		"property": "property",
	})

	tests := []struct {
		test        string
		resultError gojsonschema.ResultError
		expected    validator.Violation
	}{
		{
			test:        "InternalError",
			resultError: internalError,
			expected: validator.Violation{
				Type:              0,
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
			expected: validator.Violation{
				Type:              validator.INVALID_TYPE,
				Message:           "invalid type, expected expected, actual given",
				StructuredMessage: "invalid type",
				Arguments: []any{
					"expected", "expected",
					"actual", "given",
				},
				Path:     path.Path(""),
				Property: "",
			},
		},
		{
			test:        "RequiredError",
			resultError: requiredError,
			expected: validator.Violation{
				Type:              validator.REQUIRED,
				Message:           "missing property property",
				StructuredMessage: "missing property",
				Arguments:         []any(nil),
				Path:              path.Path(""),
				Property:          "property",
			},
		},
		{
			test:        "AdditionalPropertyNotAllowedError",
			resultError: additionalPropertyNotAllowedError,
			expected: validator.Violation{
				Type:              validator.ADDITIONAL_PROPERTY_NOT_ALLOWED,
				Message:           "additional property property is not allowed",
				StructuredMessage: "additional property is not allowed",
				Arguments:         []any(nil),
				Path:              path.Path("property"),
				Property:          "",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			validator := &Validator{}

			violation := validator.violation(test.resultError)

			s.Equal(test.expected, violation)
		})
	}
}

func (s *Suite) TestValidatorPath() {
	tests := []struct {
		test     string
		field    string
		expected string
	}{
		{
			test:     "Root",
			field:    "(root)",
			expected: "",
		},
		{
			test:     "FirstLevel",
			field:    "foo",
			expected: "foo",
		},
		{
			test:     "Levels",
			field:    "foo.bar",
			expected: "foo.bar",
		},
		{
			test:     "Index",
			field:    "foo.0.bar",
			expected: "foo[0].bar",
		},
		{
			test:     "IndexLast",
			field:    "foo.0",
			expected: "foo[0]",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			validator := &Validator{}

			path := validator.path(test.field)

			s.Equal(test.expected, path.String())

		})
	}
}
