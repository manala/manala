package validator

import "strconv"

func (s *Suite) TestCompareViolations() {
	tests := []struct {
		aColumn  int
		aLine    int
		bColumn  int
		bLine    int
		expected int
	}{
		{1, 1, 1, 1, 0},
		{1, 1, 1, 2, -1},
		{1, 1, 2, 1, -1},
		{1, 1, 2, 2, -1},

		{1, 2, 1, 1, 1},
		{1, 2, 1, 2, 0},
		{1, 2, 2, 1, 1},
		{1, 2, 2, 2, -1},

		{2, 1, 1, 1, 1},
		{2, 1, 1, 2, -1},
		{2, 1, 2, 1, 0},
		{2, 1, 2, 2, -1},

		{2, 2, 1, 1, 1},
		{2, 2, 1, 2, 1},
		{2, 2, 2, 1, 1},
		{2, 2, 2, 2, 0},
	}

	for i, test := range tests {
		s.Run(strconv.Itoa(i), func() {
			a := NewViolation("a")
			a.Column = test.aColumn
			a.Line = test.aLine

			b := NewViolation("b")
			b.Column = test.bColumn
			b.Line = test.bLine

			compare := CompareViolations(a, b)

			s.Equal(test.expected, compare)
		})
	}
}
