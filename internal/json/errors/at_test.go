package errors_test

import (
	"errors"
	"testing"

	jsonerrors "github.com/manala/manala/internal/json/errors"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/suite"
)

type AtSuite struct{ suite.Suite }

func TestAtSuite(t *testing.T) {
	suite.Run(t, new(AtSuite))
}

func (s *AtSuite) Test() {
	tests := []struct {
		test     string
		src      string
		offset   int64
		expected expectation.ErrorExpectation
	}{
		{
			test:   "EmptySource",
			src:    "",
			offset: 0,
			expected: jsonerrors.Expectation{
				Position: [2]int{0, 0},
				Err:      expectation.ErrorMessage("error"),
			},
		},
		{
			test:   "Beginning",
			src:    "foo",
			offset: 1,
			expected: jsonerrors.Expectation{
				Position: [2]int{1, 1},
				Err:      expectation.ErrorMessage("error"),
			},
		},
		{
			test:   "Middle",
			src:    "foo",
			offset: 2,
			expected: jsonerrors.Expectation{
				Position: [2]int{1, 2},
				Err:      expectation.ErrorMessage("error"),
			},
		},
		{
			test:   "AfterLine",
			src:    "foo\nbar",
			offset: 5,
			expected: jsonerrors.Expectation{
				Position: [2]int{2, 1},
				Err:      expectation.ErrorMessage("error"),
			},
		},
		{
			test:   "MultipleLines",
			src:    "foo\nbar\nbaz",
			offset: 10,
			expected: jsonerrors.Expectation{
				Position: [2]int{3, 2},
				Err:      expectation.ErrorMessage("error"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := jsonerrors.At(
				errors.New("error"),
				test.src, test.offset,
			)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
