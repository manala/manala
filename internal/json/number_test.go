package json_test

import (
	gojson "encoding/json"
	"manala/internal/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NumberSuite struct{ suite.Suite }

func TestNumberSuite(t *testing.T) {
	suite.Run(t, new(NumberSuite))
}

func (s *NumberSuite) TestType() {
	tests := []struct {
		test     string
		value    any
		expected bool
	}{
		{
			test:     "Nil",
			value:    nil,
			expected: false,
		},
		{
			test:     "Number",
			value:    gojson.Number("0"),
			expected: true,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			number, ok := json.NumberType(test.value)

			s.IsType(json.Number{}, number)
			s.Equal(test.expected, ok)
		})
	}
}

func (s *NumberSuite) Test() {
	tests := []struct {
		test              string
		number            gojson.Number
		expectedInt       int
		expectedNormalize any
	}{
		{
			test:              "Zero",
			number:            gojson.Number("0"),
			expectedInt:       0,
			expectedNormalize: 0,
		},
		{
			test:              "Integer",
			number:            gojson.Number("3"),
			expectedInt:       3,
			expectedNormalize: int64(3),
		},
		{
			test:              "FloatZero",
			number:            gojson.Number("3.0"),
			expectedInt:       0,
			expectedNormalize: 3.0,
		},
		{
			test:              "FloatMidLow",
			number:            gojson.Number("3.25"),
			expectedInt:       0,
			expectedNormalize: 3.25,
		},
		{
			test:              "FloatMid",
			number:            gojson.Number("3.5"),
			expectedInt:       0,
			expectedNormalize: 3.5,
		},
		{
			test:              "FloatMidHigh",
			number:            gojson.Number("3.75"),
			expectedInt:       0,
			expectedNormalize: 3.75,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			number := json.Number{Number: test.number}

			s.Equal(test.expectedInt, number.Int())
		})
	}
}
