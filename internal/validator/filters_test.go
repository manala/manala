package validator

import (
	"regexp"
)

func (s *Suite) TestFilter() {
	tests := []struct {
		test                      string
		filter                    Filter
		expectedMessage           string
		expectedStructuredMessage string
	}{
		{
			test: "PathMatch",
			filter: Filter{
				Path:              "path",
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "message",
			expectedStructuredMessage: "structured message",
		},
		{
			test: "PathNotMatch",
			filter: Filter{
				Path:              "foo",
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "",
			expectedStructuredMessage: "",
		},
		{
			test: "PathRegexMatch",
			filter: Filter{
				PathRegex:         regexp.MustCompile(`^path$`),
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "message",
			expectedStructuredMessage: "structured message",
		},
		{
			test: "PathRegexNotMatch",
			filter: Filter{
				PathRegex:         regexp.MustCompile(`^foo$`),
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "",
			expectedStructuredMessage: "",
		},
		{
			test: "TypeMatch",
			filter: Filter{
				Type:              REQUIRED,
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "message",
			expectedStructuredMessage: "structured message",
		},
		{
			test: "TypeNotMatch",
			filter: Filter{
				Type:              INVALID_TYPE,
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "",
			expectedStructuredMessage: "",
		},
		{
			test: "PropertyMatch",
			filter: Filter{
				Property:          "property",
				Message:           "message",
				StructuredMessage: "structured message",
			},
			expectedMessage:           "message",
			expectedStructuredMessage: "structured message",
		},
		{
			test: "PropertyNotMatch",
			filter: Filter{
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
			violation := NewViolation("")
			violation.Path = "path"
			violation.Type = REQUIRED
			violation.Property = "property"

			test.filter.Format(&violation)

			s.Equal(test.expectedMessage, violation.Message)
			s.Equal(test.expectedStructuredMessage, violation.StructuredMessage)
		})
	}
}
