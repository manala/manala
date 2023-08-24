package template

import (
	"github.com/stretchr/testify/suite"
	"manala/internal/errors/serrors"
	"testing"
	"text/template"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestError() {
	tests := []struct {
		test     string
		err      error
		expected *serrors.Assert
	}{
		{
			test: "Unknown",
			err:  serrors.New(`error`),
			expected: &serrors.Assert{
				Type:    &Error{},
				Message: "error",
			},
		},
		{
			test: "TemplateLine",
			err:  serrors.New(`template: foo.tmpl:123: message`),
			expected: &serrors.Assert{
				Type:    &Error{},
				Message: "message",
				Arguments: []any{
					"line", 123,
				},
			},
		},
		{
			test: "ContentLineColumn",
			err:  template.ExecError{Err: serrors.New(`template: :1:3: executing "" at <.foo>: nil data; no entry for key "foo"`)},
			expected: &serrors.Assert{
				Type:    &Error{},
				Message: "nil data; no entry for key \"foo\"",
				Arguments: []any{
					"context", ".foo",
					"line", 1,
					"column", 3,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := NewError(test.err)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}
