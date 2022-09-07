package template

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	"testing"
	textTemplate "text/template"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestParsingError() {
	tests := []struct {
		name   string
		err    error
		report *internalReport.Assert
	}{
		{
			name: "Unknown",
			err:  fmt.Errorf(`error`),
			report: &internalReport.Assert{
				Err: "error",
			},
		},
		{
			name: "Template Line",
			err:  fmt.Errorf(`template: foo.tmpl:123: message`),
			report: &internalReport.Assert{
				Err: "message",
				Fields: map[string]interface{}{
					"line": 123,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := NewParsingError(test.err)

			var _parsingError *ParsingError
			s.ErrorAs(err, &_parsingError)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *ErrorsSuite) TestExecutionError() {
	tests := []struct {
		name   string
		err    error
		report *internalReport.Assert
	}{
		{
			name: "Unknown",
			err:  fmt.Errorf(`error`),
			report: &internalReport.Assert{
				Err: "error",
			},
		},
		{
			name: "Content Line Column",
			err:  textTemplate.ExecError{Err: fmt.Errorf(`template: :1:3: executing "" at <.foo>: nil data; no entry for key "foo"`)},
			report: &internalReport.Assert{
				Err: `nil data; no entry for key "foo"`,
				Fields: map[string]interface{}{
					"line":    1,
					"column":  3,
					"context": ".foo",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := NewExecutionError(test.err)

			var _executionError *ExecutionError
			s.ErrorAs(err, &_executionError)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}
