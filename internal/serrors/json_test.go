package serrors_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/manala/manala/internal/serrors"

	"github.com/stretchr/testify/suite"
)

type JSONSuite struct{ suite.Suite }

func TestJSONSuite(t *testing.T) {
	suite.Run(t, new(JSONSuite))
}

func (s *JSONSuite) Test() {
	tests := []struct {
		test     string
		err      error
		expected *serrors.Assertion
	}{
		{
			test: "Unknown",
			err:  serrors.New("unknown"),
			expected: &serrors.Assertion{
				Message: "unknown",
			},
		},
		{
			test: "SyntaxError",
			err:  &json.SyntaxError{Offset: 123},
			expected: &serrors.Assertion{
				Message: "",
				Arguments: []any{
					"offset", int64(123),
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
				Type:   reflect.TypeOf(0.0),
			},
			expected: &serrors.Assertion{
				Message: "cannot unmarshal into struct field",
				Arguments: []any{
					"offset", int64(123),
					"value", "value",
					"struct", "struct",
					"field", "field",
					"type", "float64",
				},
			},
		},
		{
			test: "UnmarshalTypeError",
			err: &json.UnmarshalTypeError{
				Offset: 123,
				Value:  "value",
				Type:   reflect.TypeOf(0.0),
			},
			expected: &serrors.Assertion{
				Message: "cannot unmarshal into value",
				Arguments: []any{
					"offset", int64(123),
					"value", "value",
					"type", "float64",
				},
			},
		},
		{
			test: "InvalidUnmarshalError",
			err: &json.InvalidUnmarshalError{
				Type: reflect.TypeOf(0.0),
			},
			expected: &serrors.Assertion{
				Message: "invalid unmarshal argument",
				Arguments: []any{
					"type", "float64",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := serrors.NewJSON(test.err)

			serrors.Equal(s.T(), test.expected, err)
		})
	}
}
