package number_test

import (
	"encoding/json"
	"testing"

	"github.com/manala/manala/internal/json/number"

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
			value:    json.Number("0"),
			expected: true,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			num, ok := number.NumberType(test.value)

			s.Require().IsType(number.Number{}, num)
			s.Equal(test.expected, ok)
		})
	}
}

func (s *NumberSuite) Test() {
	tests := []struct {
		test              string
		number            json.Number
		expectedInt       int
		expectedNormalize any
	}{
		{
			test:              "Zero",
			number:            json.Number("0"),
			expectedInt:       0,
			expectedNormalize: 0,
		},
		{
			test:              "Integer",
			number:            json.Number("3"),
			expectedInt:       3,
			expectedNormalize: int64(3),
		},
		{
			test:              "FloatZero",
			number:            json.Number("3.0"),
			expectedInt:       0,
			expectedNormalize: 3.0,
		},
		{
			test:              "FloatMidLow",
			number:            json.Number("3.25"),
			expectedInt:       0,
			expectedNormalize: 3.25,
		},
		{
			test:              "FloatMid",
			number:            json.Number("3.5"),
			expectedInt:       0,
			expectedNormalize: 3.5,
		},
		{
			test:              "FloatMidHigh",
			number:            json.Number("3.75"),
			expectedInt:       0,
			expectedNormalize: 3.75,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			num := number.Number{Number: test.number}

			s.Equal(test.expectedInt, num.Int())
		})
	}
}
