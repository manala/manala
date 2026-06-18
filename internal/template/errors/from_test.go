package errors_test

import (
	"errors"
	"strings"
	"testing"
	"text/template"

	templateerrors "github.com/manala/manala/internal/template/errors"
	templateerrorstest "github.com/manala/manala/internal/template/errors/errorstest"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/suite"
)

type FromSuite struct{ suite.Suite }

func TestFromSuite(t *testing.T) {
	suite.Run(t, new(FromSuite))
}

func (s *FromSuite) Test() {
	// Source with 11 empty lines + line 12 of 40 ASCII chars (byte offset == rune offset)
	asciiSrc := strings.Repeat("\n", 11) + strings.Repeat("a", 40)

	tests := []struct {
		test     string
		err      error
		src      string
		expected expectation.ErrorExpectation
	}{
		{
			test:     "Unknown",
			err:      errors.New("unknown"),
			expected: expectation.ErrorEqual(errors.New("unknown")),
		},
		{
			test: "TextTemplateLine",
			err:  errors.New(`template: foo.tmpl:12: message`),
			expected: templateerrorstest.Expectation{
				Position: [2]int{12, 0},
				Err:      expectation.ErrorMessage("message"),
			},
		},
		{
			test: "TextEmptyTemplateLine",
			err:  errors.New(`template: :12: message`),
			expected: templateerrorstest.Expectation{
				Position: [2]int{12, 0},
				Err:      expectation.ErrorMessage("message"),
			},
		},
		{
			test: "TextTemplateLineColumnContext",
			err:  template.ExecError{Err: errors.New(`template: foo.gohtml:12:34: executing "title" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			src:  asciiSrc,
			expected: templateerrorstest.Expectation{
				Position: [2]int{12, 35},
				Err:      expectation.ErrorMessage("can't evaluate field Bar in type []foo.Bar"),
			},
		},
		{
			test: "TextEmptyTemplateLineColumnContext",
			err:  template.ExecError{Err: errors.New(`template: :12:34: executing "title" at <.Bar>: can't evaluate field Bar in type []foo.Bar`)},
			src:  asciiSrc,
			expected: templateerrorstest.Expectation{
				Position: [2]int{12, 35},
				Err:      expectation.ErrorMessage("can't evaluate field Bar in type []foo.Bar"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := templateerrors.From(test.err, test.src)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
