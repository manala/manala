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
					"template", "foo.tmpl",
					"line", 123,
				},
			},
		},
		{
			test: "ContextLineColumn",
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
		{
			test: "ContextTemplateLineColumn",
			err:  template.ExecError{Err: serrors.New(`template: message.gohtml:3:23: executing "title$htmltemplate_stateRCDATA_elementTitle" at <.Message>: can't evaluate field Message in type []app.Recipe`)},
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "can't evaluate field Message in type []app.Recipe",
				Arguments: []any{
					"context", ".Message",
					"template", "message.gohtml",
					"line", 3,
					"column", 23,
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
