package engine_test

import (
	"errors"
	"strings"
	"testing"
	textTemplate "text/template"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/template/engine"
	"github.com/manala/manala/internal/testing/expect"

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
		expected expect.ErrorExpectation
	}{
		{
			test: "Unknown",
			err:  errors.New("unknown"),
			expected: parsing.ErrorExpectation{
				Err: expect.ErrorMessageExpectation("unknown"),
			},
		},
		{
			test: "TextTemplateLine",
			err:  errors.New(`template: foo.tmpl:12: message`),
			expected: parsing.ErrorExpectation{
				Line: 12,
				Err:  expect.ErrorMessageExpectation("message"),
			},
		},
		{
			test: "TextEmptyTemplateLine",
			err:  errors.New(`template: :12: message`),
			expected: parsing.ErrorExpectation{
				Line: 12,
				Err:  expect.ErrorMessageExpectation("message"),
			},
		},
		{
			test: "TextTemplateLineColumnContext",
			err:  textTemplate.ExecError{Err: errors.New(`template: foo.gohtml:12:34: executing "title" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			src:  asciiSrc,
			expected: parsing.ErrorExpectation{
				Line:   12,
				Column: 35,
				Err:    expect.ErrorMessageExpectation("can't evaluate field Bar in type []foo.Bar"),
			},
		},
		{
			test: "TextEmptyTemplateLineColumnContext",
			err:  textTemplate.ExecError{Err: errors.New(`template: :12:34: executing "title" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			src:  asciiSrc,
			expected: parsing.ErrorExpectation{
				Line:   12,
				Column: 35,
				Err:    expect.ErrorMessageExpectation("can't evaluate field Bar in type []foo.Bar"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := engine.ErrorFrom(test.err, test.src)

			expect.Error(s.T(), test.expected, err)
		})
	}
}
