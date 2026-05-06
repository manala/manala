package errors_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	jsonerrors "github.com/manala/manala/internal/json/errors"
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
		expected expectation.ErrorExpectation
	}{
		{
			test:     "Unknown",
			err:      errors.New("unknown"),
			expected: expectation.ErrorEqual(errors.New("unknown")),
		},
		{
			test: "SyntaxError",
			err:  &json.SyntaxError{Offset: 0},
			expected: jsonerrors.Expectation{
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
			expected: jsonerrors.Expectation{
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
			expected: jsonerrors.Expectation{
				Position: [2]int{0, 0},
				Err:      expectation.ErrorMessage("wrong foo value type"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := jsonerrors.From(test.err, "")

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
