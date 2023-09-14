package json

import (
	"encoding/json"
	"manala/internal/serrors"
	"reflect"
)

func (s *Suite) TestError() {
	tests := []struct {
		test     string
		err      error
		expected *serrors.Assert
	}{
		{
			test: "Unknown",
			err:  serrors.New("unknown"),
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "unknown",
			},
		},
		{
			test: "SyntaxError",
			err:  &json.SyntaxError{Offset: 123},
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid unmarshal argument",
				Arguments: []any{
					"type", "float64",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := NewError(test.err)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}
