package serrors

import (
	"encoding/json"
	"reflect"
)

func (s *Suite) TestJson() {
	tests := []struct {
		test     string
		err      error
		expected *Assertion
	}{
		{
			test: "Unknown",
			err:  New("unknown"),
			expected: &Assertion{
				Message: "unknown",
			},
		},
		{
			test: "SyntaxError",
			err:  &json.SyntaxError{Offset: 123},
			expected: &Assertion{
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
			expected: &Assertion{
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
			expected: &Assertion{
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
			expected: &Assertion{
				Message: "invalid unmarshal argument",
				Arguments: []any{
					"type", "float64",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := NewJSON(test.err)

			Equal(s.T(), test.expected, err)
		})
	}
}
