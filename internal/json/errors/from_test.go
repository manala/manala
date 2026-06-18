package errors_test

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"

	jsonerrors "github.com/manala/manala/internal/json/errors"
	jsonerrorstest "github.com/manala/manala/internal/json/errors/errorstest"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/suite"
)

type FromSuite struct{ suite.Suite }

func TestFromSuite(t *testing.T) {
	suite.Run(t, new(FromSuite))
}

func (s *FromSuite) Test() {
	tests := []struct {
		test     string
		err      error
		src      string
		expected expectation.ErrorExpectation
	}{
		{
			test:     "Unknown",
			err:      errors.New("unknown"),
			src:      "",
			expected: expectation.ErrorEqual(errors.New("unknown")),
		},
		{
			test: "SyntaxError",
			err:  &json.SyntaxError{Offset: 0},
			src:  "",
			expected: jsonerrorstest.Expectation{
				Position: [2]int{0, 0},
				Err:      expectation.ErrorMessage(""),
			},
		},
		{
			test: "UnmarshalTypeErrorStructField",
			err: &json.UnmarshalTypeError{
				Offset: 123,
				Value:  "value",
				Struct: "struct",
				Field:  "foo",
				Type:   reflect.TypeFor[float64](),
			},
			src: "",
			expected: jsonerrorstest.Expectation{
				Position: [2]int{0, 0},
				Err:      expectation.ErrorMessage("wrong value type for field \"foo\""),
			},
		},
		{
			test: "UnmarshalTypeError",
			err: &json.UnmarshalTypeError{
				Offset: 123,
				Value:  "foo",
				Type:   reflect.TypeFor[float64](),
			},
			src: "",
			expected: jsonerrorstest.Expectation{
				Position: [2]int{0, 0},
				Err:      expectation.ErrorMessage("wrong foo value type"),
			},
		},
		{
			test: "EOFEmpty",
			err:  io.EOF,
			src:  "",
			expected: jsonerrorstest.Expectation{
				Position: [2]int{0, 0},
				Err:      expectation.ErrorMessage("EOF"),
			},
		},
		{
			test: "EOF",
			err:  io.EOF,
			src:  "{\"foo\":",
			expected: jsonerrorstest.Expectation{
				Position: [2]int{1, 7},
				Err:      expectation.ErrorMessage("EOF"),
			},
		},
		{
			test: "UnexpectedEOF",
			err:  io.ErrUnexpectedEOF,
			src:  "{\"foo\":",
			expected: jsonerrorstest.Expectation{
				Position: [2]int{1, 7},
				Err:      expectation.ErrorMessage("unexpected EOF"),
			},
		},
		{
			test: "UnexpectedEOFMultiLine",
			err:  io.ErrUnexpectedEOF,
			src:  "{\n  \"foo\":",
			expected: jsonerrorstest.Expectation{
				Position: [2]int{2, 8},
				Err:      expectation.ErrorMessage("unexpected EOF"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := jsonerrors.From(test.err, test.src)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
