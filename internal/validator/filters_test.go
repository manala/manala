package validator_test

import (
	"regexp"
	"testing"

	"manala/internal/validator"

	"github.com/stretchr/testify/suite"
)

type FiltersSuite struct{ suite.Suite }

func TestFiltersSuite(t *testing.T) {
	suite.Run(t, new(FiltersSuite))
}

func (s *FiltersSuite) Test() {
	tests := []struct {
		test                      string
		filter                    validator.Filter
		expectedMessage           string
		expectedStructuredMessage string
	}{
		{
			test: "PathMatch",
			filter: validator.Filter{
				Path:              "path",
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "message",
			expectedStructuredMessage: "structured message",
		},
		{
			test: "PathNotMatch",
			filter: validator.Filter{
				Path:              "foo",
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "",
			expectedStructuredMessage: "",
		},
		{
			test: "PathRegexMatch",
			filter: validator.Filter{
				PathRegex:         regexp.MustCompile(`^path$`),
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "message",
			expectedStructuredMessage: "structured message",
		},
		{
			test: "PathRegexNotMatch",
			filter: validator.Filter{
				PathRegex:         regexp.MustCompile(`^foo$`),
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "",
			expectedStructuredMessage: "",
		},
		{
			test: "TypeMatch",
			filter: validator.Filter{
				Type:              validator.Required,
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "message",
			expectedStructuredMessage: "structured message",
		},
		{
			test: "TypeNotMatch",
			filter: validator.Filter{
				Type:              validator.InvalidType,
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "",
			expectedStructuredMessage: "",
		},
		{
			test: "PropertyMatch",
			filter: validator.Filter{
				Property:          "property",
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "message",
			expectedStructuredMessage: "structured message",
		},
		{
			test: "PropertyNotMatch",
			filter: validator.Filter{
				Property:          "foo",
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "",
			expectedStructuredMessage: "",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			violation := validator.NewViolation("")
			violation.Path = "path"
			violation.Type = validator.Required
			violation.Property = "property"

			test.filter.Format(&violation)

			s.Equal(test.expectedMessage, violation.Message)
			s.Equal(test.expectedStructuredMessage, violation.StructuredMessage)
		})
	}
}
