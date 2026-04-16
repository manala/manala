package unmarshaler_test

import (
	"encoding/json"
	"testing"

	"github.com/manala/manala/internal/json/unmarshaler"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/testing/expect"

	"github.com/stretchr/testify/suite"
)

type UnmarshalSuite struct{ suite.Suite }

func TestUnmarshalSuite(t *testing.T) {
	suite.Run(t, new(UnmarshalSuite))
}

func (s *UnmarshalSuite) TestErrors() {
	tests := []struct {
		test     string
		data     string
		expected expect.ErrorExpectation
	}{
		{
			test: "Syntax",
			data: `foo`,
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 2,
				Err:    expect.ErrorMessageExpectation("invalid character 'o' in literal false (expecting 'a')"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := unmarshaler.Unmarshal([]byte(test.data), nil)

			expect.Error(s.T(), test.expected, err)
		})
	}
}

func (s *UnmarshalSuite) Test() {
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

			err := unmarshaler.Unmarshal([]byte(test.data), &value)

			s.Require().NoError(err)
			s.Equal(test.expected, value)
		})
	}
}
