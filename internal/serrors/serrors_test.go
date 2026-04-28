package serrors_test

import (
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/expect"

	"github.com/stretchr/testify/suite"
)

type SerrorsSuite struct{ suite.Suite }

func TestSerrorsSuite(t *testing.T) {
	suite.Run(t, new(SerrorsSuite))
}

func (s *SerrorsSuite) TestError() {
	s.Run("New", func() {
		err := serrors.New("error")

		expect.Error(s.T(), serrors.Expectation{
			Message: "error",
		}, err)
	})

	s.Run("Arguments", func() {
		foo := "foo"
		bar := "bar"

		err := serrors.New("error").
			With(foo, bar)

		expect.Error(s.T(), serrors.Expectation{
			Message: "error",
			Attrs: [][2]any{
				{foo, bar},
			},
		}, err)
	})

	s.Run("Dump", func() {
		dump := "dump"

		err := serrors.New("error").
			WithDump(dump)

		expect.Error(s.T(), serrors.Expectation{
			Message: "error",
			Dump:    dump,
		}, err)
	})

	s.Run("Dump", func() {
		dump := "dump"

		err := serrors.New("error").
			WithDump(dump)

		expect.Error(s.T(), serrors.Expectation{
			Message: "error",
			Dump:    dump,
		}, err)
	})

	s.Run("Errors", func() {
		foo := serrors.New("foo")
		bar := serrors.New("bar")

		err := serrors.New("error").
			WithErrors(foo, bar)

		expect.Error(s.T(), serrors.Expectation{
			Message: "error",
			Errors: []expect.ErrorExpectation{
				serrors.Expectation{
					Message: "foo",
				},
				serrors.Expectation{
					Message: "bar",
				},
			},
		}, err)
	})

	s.Run("NilErrors", func() {
		foo := serrors.New("foo")

		err := serrors.New("error").
			WithErrors(nil, foo, nil)

		expect.Error(s.T(), serrors.Expectation{
			Message: "error",
			Errors: []expect.ErrorExpectation{
				serrors.Expectation{
					Message: "foo",
				},
			},
		}, err)
	})
}
