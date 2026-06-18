package decoder_test

import (
	"encoding/json"
	"testing"

	jsondecoder "github.com/manala/manala/internal/json/decoder"
	jsonerrorstest "github.com/manala/manala/internal/json/errors/errorstest"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/suite"
)

type DecodeSuite struct{ suite.Suite }

func TestDecodeSuite(t *testing.T) {
	suite.Run(t, new(DecodeSuite))
}

func (s *DecodeSuite) TestErrors() {
	tests := []struct {
		test     string
		data     string
		expected expectation.ErrorExpectation
	}{
		{
			test: "Syntax",
			data: `foo`,
			expected: jsonerrorstest.Expectation{
				Position: [2]int{1, 2},
				Err:      expectation.ErrorMessage("invalid character 'o' in literal false (expecting 'a')"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := jsondecoder.Decode([]byte(test.data), nil)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}

func (s *DecodeSuite) Test() {
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

			err := jsondecoder.Decode([]byte(test.data), &value)

			s.Require().NoError(err)
			s.Equal(test.expected, value)
		})
	}
}
