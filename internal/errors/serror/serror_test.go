package serror_test

import (
	"errors"
	"testing"

	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/suite"
)

type SerrorSuite struct{ suite.Suite }

func TestSerrorSuite(t *testing.T) {
	suite.Run(t, new(SerrorSuite))
}

func (s *SerrorSuite) TestError() {
	s.Run("New", func() {
		err := serror.New("error")

		expectation.ExpectError(s.T(), serror.Expectation{
			Msg: "error",
		}, err)
	})

	s.Run("Attrs", func() {
		foo := "foo"
		bar := "bar"

		err := serror.New("error").
			With(foo, bar)

		expectation.ExpectError(s.T(), serror.Expectation{
			Msg: "error",
			Attrs: [][2]any{
				{foo, bar},
			},
		}, err)
	})

	s.Run("Dump", func() {
		dump := "dump"

		err := serror.New("error").
			WithDump(dump)

		expectation.ExpectError(s.T(), serror.Expectation{
			Msg:  "error",
			Dump: dump,
		}, err)
	})

	s.Run("Err", func() {
		foo := errors.New("foo")

		err := serror.New("error").
			WithErr(foo)

		expectation.ExpectError(s.T(), serror.Expectation{
			Msg: "error",
			Err: expectation.ErrorEqual(foo),
		}, err)
	})
}
