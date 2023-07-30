package validation

import (
	"github.com/stretchr/testify/suite"
	"github.com/xeipuuv/gojsonschema"
	"manala/internal/errors/serrors"
	"regexp"
	"testing"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestError() {
	result := &gojsonschema.Result{}

	resultError := &gojsonschema.InternalError{}
	resultError.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))

	result.AddError(resultError, gojsonschema.ErrorDetails{})

	err := NewError("error", result)

	serrors.Equal(s.Assert(), &serrors.Assert{
		Type:    &Error{},
		Message: "error",
		Errors: []*serrors.Assert{
			{
				Type:    &ResultError{},
				Message: "(root): ",
			},
		},
	}, err)
}

func (s *ErrorsSuite) TestResultError() {
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
		test     string
		err      gojsonschema.ResultError
		expected *serrors.Assert
	}{
		{
			test: "InvalidTypeError",
			err:  invalidTypeError,
			expected: &serrors.Assert{
				Type:    &ResultError{},
				Message: "invalid type",
				Arguments: []any{
					"expected", "expected",
					"given", "given",
				},
			},
		},
		{
			test: "RequiredError",
			err:  requiredError,
			expected: &serrors.Assert{
				Type:    &ResultError{},
				Message: "missing property",
				Arguments: []any{
					"property", "property",
				},
			},
		},
		{
			test: "AdditionalPropertyNotAllowedError",
			err:  additionalPropertyNotAllowedError,
			expected: &serrors.Assert{
				Type:    &ResultError{},
				Message: "additional property is not allowed",
				Arguments: []any{
					"property", "property",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := NewResultError(test.err, []ErrorMessage{})

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *ErrorsSuite) TestErrorMessage() {
	result := &gojsonschema.ResultErrorFields{}
	result.SetContext(gojsonschema.NewJsonContext(gojsonschema.STRING_CONTEXT_ROOT, nil))
	result.SetType("type")
	result.SetDetails(gojsonschema.ErrorDetails{
		"property": "property",
	})

	tests := []struct {
		test            string
		message         ErrorMessage
		expectedMessage string
		expectedMatch   bool
	}{
		{
			test: "FieldMatch",
			message: ErrorMessage{
				Message: "message",
				Field:   "(root)",
			},
			expectedMessage: "message",
			expectedMatch:   true,
		},
		{
			test: "FieldNotMatch",
			message: ErrorMessage{
				Message: "message",
				Field:   "(foo)",
			},
			expectedMessage: "",
			expectedMatch:   false,
		},
		{
			test: "FieldRegexMatch",
			message: ErrorMessage{
				Message:    "message",
				FieldRegex: regexp.MustCompile(`\(root\)`),
			},
			expectedMessage: "message",
			expectedMatch:   true,
		},
		{
			test: "FieldRegexNotMatch",
			message: ErrorMessage{
				Message:    "message",
				FieldRegex: regexp.MustCompile(`\(foo\)`),
			},
			expectedMessage: "",
			expectedMatch:   false,
		},
		{
			test: "TypeMatch",
			message: ErrorMessage{
				Message: "message",
				Type:    "type",
			},
			expectedMessage: "message",
			expectedMatch:   true,
		},
		{
			test: "TypeNotMatch",
			message: ErrorMessage{
				Message: "message",
				Type:    "foo",
			},
			expectedMessage: "",
			expectedMatch:   false,
		},
		{
			test: "PropertyMatch",
			message: ErrorMessage{
				Message:  "message",
				Property: "property",
			},
			expectedMessage: "message",
			expectedMatch:   true,
		},
		{
			test: "PropertyNotMatch",
			message: ErrorMessage{
				Message:  "message",
				Property: "foo",
			},
			expectedMessage: "",
			expectedMatch:   false,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			message, match := test.message.Match(result)

			s.Equal(test.expectedMessage, message)
			s.Equal(test.expectedMatch, match)
		})
	}
}
