package template

import (
	"manala/internal/serrors"
	"text/template"
)

func (s *Suite) TestError() {
	tests := []struct {
		test     string
		err      error
		expected *serrors.Assert
	}{
		{
			test: "Unknown",
			err:  serrors.New(`error`),
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "error",
			},
		},
		{
			test: "TemplateLine",
			err:  serrors.New(`template: foo.tmpl:123: message`),
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
				Type:    serrors.Error{},
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
