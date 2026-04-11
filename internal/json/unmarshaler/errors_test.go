package unmarshaler_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/manala/manala/internal/json/unmarshaler"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"

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
		expected errors.Assertion
	}{
		{
			test:   "EmptySource",
			src:    "",
			offset: 0,
			expected: &parsing.ErrorAssertion{
				Line:   0,
				Column: 0,
				Err: &serrors.Assertion{
					Message: "error",
				},
			},
		},
		{
			test:   "Beginning",
			src:    "foo",
			offset: 1,
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "error",
				},
			},
		},
		{
			test:   "Middle",
			src:    "foo",
			offset: 2,
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 2,
				Err: &serrors.Assertion{
					Message: "error",
				},
			},
		},
		{
			test:   "AfterLine",
			src:    "foo\nbar",
			offset: 5,
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "error",
				},
			},
		},
		{
			test:   "MultipleLines",
			src:    "foo\nbar\nbaz",
			offset: 10,
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 2,
				Err: &serrors.Assertion{
					Message: "error",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := unmarshaler.ErrorAt(
				serrors.New("error"),
				test.src, test.offset,
			)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *ErrorsSuite) TestErrorFrom() {
	tests := []struct {
		test     string
		err      error
		expected errors.Assertion
	}{
		{
			test: "Unknown",
			err:  serrors.New("unknown"),
			expected: &parsing.ErrorAssertion{
				Line:   0,
				Column: 0,
				Err: &serrors.Assertion{
					Message: "unknown",
				},
			},
		},
		{
			test: "SyntaxError",
			err:  &json.SyntaxError{Offset: 0},
			expected: &parsing.ErrorAssertion{
				Line:   0,
				Column: 0,
				Err: &serrors.Assertion{
					Message: "",
				},
			},
		},
		{
			test: "UnmarshalTypeErrorStructField",
			err: &json.UnmarshalTypeError{
				Offset: 123,
				Value:  "value",
				Struct: "struct",
				Field:  "field",
				Type:   reflect.TypeFor[float64](),
			},
			expected: &parsing.ErrorAssertion{
				Line:   0,
				Column: 0,
				Err: &serrors.Assertion{
					Message: "cannot unmarshal into struct field",
					Arguments: []any{
						"value", "value",
						"struct", "struct",
						"field", "field",
						"type", "float64",
					},
				},
			},
		},
		{
			test: "UnmarshalTypeError",
			err: &json.UnmarshalTypeError{
				Offset: 123,
				Value:  "value",
				Type:   reflect.TypeFor[float64](),
			},
			expected: &parsing.ErrorAssertion{
				Line:   0,
				Column: 0,
				Err: &serrors.Assertion{
					Message: "cannot unmarshal into value",
					Arguments: []any{
						"value", "value",
						"type", "float64",
					},
				},
			},
		},
		{
			test: "InvalidUnmarshalError",
			err: &json.InvalidUnmarshalError{
				Type: reflect.TypeFor[float64](),
			},
			expected: &parsing.ErrorAssertion{
				Line:   0,
				Column: 0,
				Err: &serrors.Assertion{
					Message: "invalid unmarshal argument",
					Arguments: []any{
						"type", "float64",
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := unmarshaler.ErrorFrom(test.err, "")

			errors.Equal(s.T(), test.expected, err)
		})
	}
}
