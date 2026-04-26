package unmarshaler_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/manala/manala/internal/json/unmarshaler"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/testing/expect"

	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestErrorAt() {
	tests := []struct {
		test     string
		src      string
		offset   int64
		expected expect.ErrorExpectation
	}{
		{
			test:   "EmptySource",
			src:    "",
			offset: 0,
			expected: parsing.ErrorExpectation{
				Line:   0,
				Column: 0,
				Err:    expect.ErrorMessageExpectation("error"),
			},
		},
		{
			test:   "Beginning",
			src:    "foo",
			offset: 1,
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("error"),
			},
		},
		{
			test:   "Middle",
			src:    "foo",
			offset: 2,
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 2,
				Err:    expect.ErrorMessageExpectation("error"),
			},
		},
		{
			test:   "AfterLine",
			src:    "foo\nbar",
			offset: 5,
			expected: parsing.ErrorExpectation{
				Line:   2,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("error"),
			},
		},
		{
			test:   "MultipleLines",
			src:    "foo\nbar\nbaz",
			offset: 10,
			expected: parsing.ErrorExpectation{
				Line:   3,
				Column: 2,
				Err:    expect.ErrorMessageExpectation("error"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := unmarshaler.ErrorAt(
				errors.New("error"),
				test.src, test.offset,
			)

			expect.Error(s.T(), test.expected, err)
		})
	}
}

func (s *ErrorsSuite) TestErrorFrom() {
	tests := []struct {
		test     string
		err      error
		expected expect.ErrorExpectation
	}{
		{
			test: "Unknown",
			err:  errors.New("unknown"),
			expected: parsing.ErrorExpectation{
				Line:   0,
				Column: 0,
				Err:    expect.ErrorMessageExpectation("unknown"),
			},
		},
		{
			test: "SyntaxError",
			err:  &json.SyntaxError{Offset: 0},
			expected: parsing.ErrorExpectation{
				Line:   0,
				Column: 0,
				Err:    expect.ErrorMessageExpectation(""),
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
			expected: parsing.ErrorExpectation{
				Line:   0,
				Column: 0,
				Err:    expect.ErrorMessageExpectation("wrong value type for field \"foo\""),
			},
		},
		{
			test: "UnmarshalTypeError",
			err: &json.UnmarshalTypeError{
				Offset: 123,
				Value:  "foo",
				Type:   reflect.TypeFor[float64](),
			},
			expected: parsing.ErrorExpectation{
				Line:   0,
				Column: 0,
				Err:    expect.ErrorMessageExpectation("wrong foo value type"),
			},
		},
		{
			test: "InvalidUnmarshalError",
			err: &json.InvalidUnmarshalError{
				Type: reflect.TypeFor[float64](),
			},
			expected: parsing.ErrorExpectation{
				Line:   0,
				Column: 0,
				Err:    expect.ErrorMessageExpectation("invalid unmarshal argument"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := unmarshaler.ErrorFrom(test.err, "")

			expect.Error(s.T(), test.expected, err)
		})
	}
}
