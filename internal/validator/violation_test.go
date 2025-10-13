package validator_test

import (
	"strconv"
	"testing"

	"manala/internal/validator"

	"github.com/stretchr/testify/suite"
)

type ViolationSuite struct{ suite.Suite }

func TestViolationSuite(t *testing.T) {
	suite.Run(t, new(ViolationSuite))
}

func (s *ViolationSuite) TestCompare() {
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
			a := validator.NewViolation("a")
			a.Column = test.aColumn
			a.Line = test.aLine

			b := validator.NewViolation("b")
			b.Column = test.bColumn
			b.Line = test.bLine

			compare := validator.CompareViolations(a, b)

			s.Equal(test.expected, compare)
		})
	}
}
