package template

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	internalErrors "manala/internal/errors"
	"testing"
	textTemplate "text/template"
)

var internalError *internalErrors.InternalError

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestParsingError() {
	tests := []struct {
		name             string
		err              error
		wantMessage      string
		wantFieldLine    int
		wantFieldMessage string
	}{
		{
			name:        "Unknown",
			err:         fmt.Errorf(`foo`),
			wantMessage: "template error",
		},
		{
			name:             "Template Line",
			err:              fmt.Errorf(`template: foo.tmpl:123: Message`),
			wantMessage:      "template error",
			wantFieldLine:    123,
			wantFieldMessage: `Message`,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {

			// Without path
			err := ParsingError(test.err)
			s.ErrorAs(err, &internalError)
			s.Equal(test.wantMessage, internalError.Message)
			if test.wantFieldLine != 0 {
				s.Equal(test.wantFieldLine, internalError.Fields["line"])
			} else {
				s.NotContains(internalError.Fields, "line")
			}
			if test.wantFieldMessage != "" {
				s.Equal(test.wantFieldMessage, internalError.Fields["message"])
			} else {
				s.NotContains(internalError.Fields, "message")
			}

			// With path
			err = ParsingPathError("path", test.err)
			s.ErrorAs(err, &internalError)
			s.Equal(test.wantMessage, internalError.Message)
			if test.wantFieldLine != 0 {
				s.Equal(test.wantFieldLine, internalError.Fields["line"])
			} else {
				s.NotContains(internalError.Fields, "line")
			}
			if test.wantFieldMessage != "" {
				s.Equal(test.wantFieldMessage, internalError.Fields["message"])
			} else {
				s.NotContains(internalError.Fields, "message")
			}
			s.Equal("path", internalError.Fields["path"])
		})
	}
}

func (s *ErrorsSuite) TestExecutionError() {
	tests := []struct {
		name             string
		err              error
		wantMessage      string
		wantFieldLine    int
		wantFieldColumn  int
		wantFieldContext string
		wantFieldMessage string
	}{
		{
			name:        "Unknown",
			err:         fmt.Errorf(`foo`),
			wantMessage: "template error",
		},
		{
			name:             "Content Line Column",
			err:              textTemplate.ExecError{Err: fmt.Errorf(`template: :1:3: executing "" at <.foo>: nil data; no entry for key "foo"`)},
			wantMessage:      "template error",
			wantFieldLine:    1,
			wantFieldColumn:  3,
			wantFieldContext: ".foo",
			wantFieldMessage: `nil data; no entry for key "foo"`,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {

			// With path
			err := ExecutionError(test.err)
			s.ErrorAs(err, &internalError)
			s.Equal(test.wantMessage, internalError.Message)
			if test.wantFieldLine != 0 {
				s.Equal(test.wantFieldLine, internalError.Fields["line"])
			} else {
				s.NotContains(internalError.Fields, "line")
			}
			if test.wantFieldColumn != 0 {
				s.Equal(test.wantFieldColumn, internalError.Fields["column"])
			} else {
				s.NotContains(internalError.Fields, "column")
			}
			if test.wantFieldContext != "" {
				s.Equal(test.wantFieldContext, internalError.Fields["context"])
			} else {
				s.NotContains(internalError.Fields, "context")
			}
			if test.wantFieldMessage != "" {
				s.Equal(test.wantFieldMessage, internalError.Fields["message"])
			} else {
				s.NotContains(internalError.Fields, "message")
			}

			// Without path
			err = ExecutionPathError("path", test.err)
			s.ErrorAs(err, &internalError)
			s.Equal(test.wantMessage, internalError.Message)
			if test.wantFieldLine != 0 {
				s.Equal(test.wantFieldLine, internalError.Fields["line"])
			} else {
				s.NotContains(internalError.Fields, "line")
			}
			if test.wantFieldColumn != 0 {
				s.Equal(test.wantFieldColumn, internalError.Fields["column"])
			} else {
				s.NotContains(internalError.Fields, "column")
			}
			if test.wantFieldContext != "" {
				s.Equal(test.wantFieldContext, internalError.Fields["context"])
			} else {
				s.NotContains(internalError.Fields, "context")
			}
			if test.wantFieldMessage != "" {
				s.Equal(test.wantFieldMessage, internalError.Fields["message"])
			} else {
				s.NotContains(internalError.Fields, "message")
			}
			s.Equal("path", internalError.Fields["path"])
		})
	}
}
