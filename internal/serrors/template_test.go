package serrors

import (
	htmlTemplate "html/template"
	textTemplate "text/template"
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
				Message: `error`,
			},
		},
		{
			test: "TextTemplateLine",
			err:  New(`template: foo.tmpl:12: message`),
			expected: &Assertion{
				Message: `message`,
				Arguments: []any{
					"template", "foo.tmpl",
					"line", 12,
				},
			},
		},
		{
			test: "TextEmptyTemplateLine",
			err:  New(`template: :12: message`),
			expected: &Assertion{
				Message: `message`,
				Arguments: []any{
					"line", 12,
				},
			},
		},
		{
			test: "TextTemplateLineColumnContext",
			err:  textTemplate.ExecError{Err: New(`template: foo.gohtml:12:34: executing "title$htmltemplate_stateRCDATA_elementTitle" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			expected: &Assertion{
				Message: `can't evaluate field Bar in type []foo.Bar`,
				Arguments: []any{
					"context", ".Bar",
					"template", "foo.gohtml",
					"line", 12,
					"column", 34,
				},
			},
		},
		{
			test: "TextEmptyTemplateLineColumnContext",
			err:  textTemplate.ExecError{Err: New(`template: :12:34: executing "title$htmltemplate_stateRCDATA_elementTitle" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			expected: &Assertion{
				Message: `can't evaluate field Bar in type []foo.Bar`,
				Arguments: []any{
					"context", ".Bar",
					"line", 12,
					"column", 34,
				},
			},
		},
		{
			test: "HtmlTemplateLineColumn",
			err:  &htmlTemplate.Error{Name: `foo.gohtml:12:34`, Description: `no such template "bar"`},
			expected: &Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"template", "foo.gohtml",
					"line", 12,
					"column", 34,
				},
			},
		},
		{
			test: "HtmlEmptyTemplateLineColumn",
			err:  &htmlTemplate.Error{Name: `:12:34`, Description: `no such template "bar"`},
			expected: &Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"line", 12,
					"column", 34,
				},
			},
		},
		{
			test: "HtmlTemplateLine",
			err:  &htmlTemplate.Error{Name: `foo.gohtml:12`, Description: `no such template "bar"`},
			expected: &Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"template", "foo.gohtml",
					"line", 12,
				},
			},
		},
		{
			test: "HtmlEmptyTemplateLine",
			err:  &htmlTemplate.Error{Name: `:12`, Description: `no such template "bar"`},
			expected: &Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"line", 12,
				},
			},
		},
		{
			test: "HtmlTemplate",
			err:  &htmlTemplate.Error{Name: `foo.gohtml`, Description: `no such template "bar"`},
			expected: &Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"template", "foo.gohtml",
				},
			},
		},
		{
			test: "HtmlEmptyTemplate",
			err:  &htmlTemplate.Error{Name: `:12`, Description: `no such template "bar"`},
			expected: &Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"line", 12,
				},
			},
		},
		{
			test: "Html",
			err:  &htmlTemplate.Error{Description: `no such template "bar"`},
			expected: &Assertion{
				Message: `no such template "bar"`,
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
