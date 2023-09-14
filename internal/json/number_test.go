package json

import "encoding/json"

func (s *Suite) TestNumberType() {
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
			number, ok := NumberType(test.value)

			s.IsType(Number{}, number)
			s.Equal(test.expected, ok)
		})
	}
}

func (s *Suite) TestNumber() {
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
			number := Number{Number: test.number}

			s.Equal(test.expectedInt, number.Int())
		})
	}
}
