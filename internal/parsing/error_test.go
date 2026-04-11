package parsing_test

import (
	"errors"
	"testing"

	"github.com/manala/manala/internal/parsing"

	"github.com/stretchr/testify/suite"
)

type ErrorSuite struct{ suite.Suite }

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (s *ErrorSuite) TestFlatten() {
	s.Run("SingleLayer", func() {
		root := errors.New("root")
		err := &parsing.Error{Err: root, Line: 3, Column: 5}

		flat := err.Flatten()

		s.Same(err, flat)
	})

	s.Run("TwoLayers", func() {
		root := errors.New("root")
		err := &parsing.Error{
			Err:  &parsing.Error{Err: root, Line: 2, Column: 3},
			Line: 10, Column: 5,
		}

		flat := err.Flatten()

		s.Equal(root, flat.Err)
		s.Equal(11, flat.Line)
		s.Equal(7, flat.Column)
	})

	s.Run("ThreeLayers", func() {
		root := errors.New("root")
		err := &parsing.Error{
			Err: &parsing.Error{
				Err:    &parsing.Error{Err: root, Line: 2, Column: 3},
				Line:   10,
				Column: 5,
			},
			Line: 20, Column: 3,
		}

		flat := err.Flatten()

		s.Equal(root, flat.Err)
		s.Equal(30, flat.Line)
		s.Equal(9, flat.Column)
	})

	s.Run("InnerZeroLine", func() {
		root := errors.New("root")
		err := &parsing.Error{
			Err:    &parsing.Error{Err: root, Line: 0, Column: 3},
			Line:   10,
			Column: 5,
		}

		flat := err.Flatten()

		s.Equal(root, flat.Err)
		s.Equal(10, flat.Line)
		s.Equal(7, flat.Column)
	})

	s.Run("InnerZeroColumn", func() {
		root := errors.New("root")
		err := &parsing.Error{
			Err:    &parsing.Error{Err: root, Line: 2, Column: 0},
			Line:   10,
			Column: 5,
		}

		flat := err.Flatten()

		s.Equal(root, flat.Err)
		s.Equal(11, flat.Line)
		s.Equal(5, flat.Column)
	})

	s.Run("OuterZeroLine", func() {
		root := errors.New("root")
		err := &parsing.Error{
			Err:    &parsing.Error{Err: root, Line: 2, Column: 3},
			Line:   0,
			Column: 5,
		}

		flat := err.Flatten()

		s.Equal(root, flat.Err)
		s.Equal(2, flat.Line)
		s.Equal(7, flat.Column)
	})

	s.Run("OuterZeroColumn", func() {
		root := errors.New("root")
		err := &parsing.Error{
			Err:    &parsing.Error{Err: root, Line: 2, Column: 3},
			Line:   10,
			Column: 0,
		}

		flat := err.Flatten()

		s.Equal(root, flat.Err)
		s.Equal(11, flat.Line)
		s.Equal(3, flat.Column)
	})

	s.Run("BothZero", func() {
		root := errors.New("root")
		err := &parsing.Error{
			Err:    &parsing.Error{Err: root, Line: 0, Column: 0},
			Line:   0,
			Column: 0,
		}

		flat := err.Flatten()

		s.Equal(root, flat.Err)
		s.Equal(0, flat.Line)
		s.Equal(0, flat.Column)
	})

	s.Run("NilErr", func() {
		err := &parsing.Error{Err: nil, Line: 3, Column: 5}

		flat := err.Flatten()

		s.Same(err, flat)
	})
}
