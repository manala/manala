package engine_test

import (
	"strings"
	"testing"
	textTemplate "text/template"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/template/engine"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestErrorFrom() {
	// Source with 11 empty lines + line 12 of 40 ASCII chars (byte offset == rune offset)
	asciiSrc := strings.Repeat("\n", 11) + strings.Repeat("a", 40)

	tests := []struct {
		test     string
		err      error
		src      string
		expected errors.Assertion
	}{
		{
			test: "Unknown",
			err:  serrors.New("unknown"),
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "unknown",
				},
			},
		},
		{
			test: "TextTemplateLine",
			err:  serrors.New(`template: foo.tmpl:12: message`),
			expected: &parsing.ErrorAssertion{
				Line: 12,
				Err: &serrors.Assertion{
					Message: "message",
					Arguments: []any{
						"template", "foo.tmpl",
					},
				},
			},
		},
		{
			test: "TextEmptyTemplateLine",
			err:  serrors.New(`template: :12: message`),
			expected: &parsing.ErrorAssertion{
				Line: 12,
				Err: &serrors.Assertion{
					Message: "message",
				},
			},
		},
		{
			test: "TextTemplateLineColumnContext",
			err:  textTemplate.ExecError{Err: serrors.New(`template: foo.gohtml:12:34: executing "title" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			src:  asciiSrc,
			expected: &parsing.ErrorAssertion{
				Line:   12,
				Column: 35,
				Err: &serrors.Assertion{
					Message: "can't evaluate field Bar in type []foo.Bar",
					Arguments: []any{
						"context", ".Bar",
						"template", "foo.gohtml",
					},
				},
			},
		},
		{
			test: "TextEmptyTemplateLineColumnContext",
			err:  textTemplate.ExecError{Err: serrors.New(`template: :12:34: executing "title" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			src:  asciiSrc,
			expected: &parsing.ErrorAssertion{
				Line:   12,
				Column: 35,
				Err: &serrors.Assertion{
					Message: "can't evaluate field Bar in type []foo.Bar",
					Arguments: []any{
						"context", ".Bar",
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := engine.ErrorFrom(test.err, test.src)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}
