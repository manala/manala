package serrors

import (
	"text/template"
)

func (s *Suite) TestTemplate() {
	tests := []struct {
		test     string
		err      error
		expected *Assertion
	}{
		{
			test: "Unknown",
			err:  New(`error`),
			expected: &Assertion{
				Message: "error",
			},
		},
		{
			test: "TemplateLine",
			err:  New(`template: foo.tmpl:123: message`),
			expected: &Assertion{
				Message: "message",
				Arguments: []any{
					"template", "foo.tmpl",
					"line", 123,
				},
			},
		},
		{
			test: "ContextLineColumn",
			err:  template.ExecError{Err: New(`template: :1:3: executing "" at <.foo>: nil data; no entry for key "foo"`)},
			expected: &Assertion{
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
			err:  template.ExecError{Err: New(`template: message.gohtml:3:23: executing "title$htmltemplate_stateRCDATA_elementTitle" at <.Message>: can't evaluate field Message in type []app.Recipe`)},
			expected: &Assertion{
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
			err := NewTemplate(test.err)

			Equal(s.T(), test.expected, err)
		})
	}
}
