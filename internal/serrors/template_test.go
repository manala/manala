package serrors_test

import (
	htmlTemplate "html/template"
	"testing"
	textTemplate "text/template"

	"manala/internal/serrors"

	"github.com/stretchr/testify/suite"
)

type TemplateSuite struct{ suite.Suite }

func TestTemplateSuite(t *testing.T) {
	suite.Run(t, new(TemplateSuite))
}

func (s *TemplateSuite) Test() {
	tests := []struct {
		test     string
		err      error
		expected *serrors.Assertion
	}{
		{
			test: "Unknown",
			err:  serrors.New(`error`),
			expected: &serrors.Assertion{
				Message: `error`,
			},
		},
		{
			test: "TextTemplateLine",
			err:  serrors.New(`template: foo.tmpl:12: message`),
			expected: &serrors.Assertion{
				Message: `message`,
				Arguments: []any{
					"template", "foo.tmpl",
					"line", 12,
				},
			},
		},
		{
			test: "TextEmptyTemplateLine",
			err:  serrors.New(`template: :12: message`),
			expected: &serrors.Assertion{
				Message: `message`,
				Arguments: []any{
					"line", 12,
				},
			},
		},
		{
			test: "TextTemplateLineColumnContext",
			err:  textTemplate.ExecError{Err: serrors.New(`template: foo.gohtml:12:34: executing "title$htmltemplate_stateRCDATA_elementTitle" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			expected: &serrors.Assertion{
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
			err:  textTemplate.ExecError{Err: serrors.New(`template: :12:34: executing "title$htmltemplate_stateRCDATA_elementTitle" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			expected: &serrors.Assertion{
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
			expected: &serrors.Assertion{
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
			expected: &serrors.Assertion{
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
			expected: &serrors.Assertion{
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
			expected: &serrors.Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"line", 12,
				},
			},
		},
		{
			test: "HtmlTemplate",
			err:  &htmlTemplate.Error{Name: `foo.gohtml`, Description: `no such template "bar"`},
			expected: &serrors.Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"template", "foo.gohtml",
				},
			},
		},
		{
			test: "HtmlEmptyTemplate",
			err:  &htmlTemplate.Error{Name: `:12`, Description: `no such template "bar"`},
			expected: &serrors.Assertion{
				Message: `no such template "bar"`,
				Arguments: []any{
					"line", 12,
				},
			},
		},
		{
			test: "Html",
			err:  &htmlTemplate.Error{Description: `no such template "bar"`},
			expected: &serrors.Assertion{
				Message: `no such template "bar"`,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := serrors.NewTemplate(test.err)

			serrors.Equal(s.T(), test.expected, err)
		})
	}
}
