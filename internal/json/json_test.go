package json

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"manala/internal/serrors"
	"testing"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestUnmarshalErrors() {
	tests := []struct {
		test     string
		data     string
		expected *serrors.Assert
	}{
		{
			test: "Syntax",
			data: `foo`,
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid character 'o' in literal false (expecting 'a')",
				Arguments: []any{
					"offset", int64(2),
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := Unmarshal([]byte(test.data), nil)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestUnmarshal() {
	tests := []struct {
		test     string
		data     string
		expected map[string]any
	}{
		{
			test:     "Empty",
			data:     `{}`,
			expected: map[string]any{},
		},
		{
			test:     "Null",
			data:     `{"null": null}`,
			expected: map[string]any{"null": nil},
		},
		{
			test:     "Boolean",
			data:     `{"true": true, "false": false}`,
			expected: map[string]any{"true": true, "false": false},
		},
		{
			test:     "String",
			data:     `{"string": "string"}`,
			expected: map[string]any{"string": "string"},
		},
		{
			test:     "Integer",
			data:     `{"integer": 12}`,
			expected: map[string]any{"integer": json.Number("12")},
		},
		{
			test:     "Float",
			data:     `{"float": 2.3}`,
			expected: map[string]any{"float": json.Number("2.3")},
		},
		{
			test:     "FloatAsInteger",
			data:     `{"float": 3.0}`,
			expected: map[string]any{"float": json.Number("3.0")},
		},
		{
			test:     "FloatAsString",
			data:     `{"float": "3.0"}`,
			expected: map[string]any{"float": "3.0"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			var value map[string]any

			err := Unmarshal([]byte(test.data), &value)

			s.NoError(err)
			s.Equal(test.expected, value)
		})
	}
}
