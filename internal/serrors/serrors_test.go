package serrors_test

import (
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type SerrorsSuite struct{ suite.Suite }

func TestSerrorsSuite(t *testing.T) {
	suite.Run(t, new(SerrorsSuite))
}

func (s *SerrorsSuite) TestError() {
	s.Run("New", func() {
		err := serrors.New("error")

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "error",
		}, err)
	})

	s.Run("Message", func() {
		message := "message"

		err := serrors.New("error").
			WithMessage(message)

		errors.Equal(s.T(), &serrors.Assertion{
			Message: message,
		}, err)
	})

	s.Run("Arguments", func() {
		foo := "foo"
		bar := "bar"

		err := serrors.New("error").
			WithArguments(foo, bar)

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "error",
			Arguments: []any{
				foo, bar,
			},
		}, err)
	})

	s.Run("Dump", func() {
		dump := "dump"

		err := serrors.New("error").
			WithDump(dump)

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "error",
			Dump:    dump,
		}, err)
	})

	s.Run("Dumper", func() {
		dump := "dump"

		err := serrors.New("error").
			WithDumper(serrors.StringDumper(dump))

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "error",
			Dump:    dump,
		}, err)
	})

	s.Run("Errors", func() {
		foo := serrors.New("foo")
		bar := serrors.New("bar")

		err := serrors.New("error").
			WithErrors(foo, bar)

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "error",
			Errors: []errors.Assertion{
				&serrors.Assertion{
					Message: "foo",
				},
				&serrors.Assertion{
					Message: "bar",
				},
			},
		}, err)
	})

	s.Run("NilErrors", func() {
		foo := serrors.New("foo")

		err := serrors.New("error").
			WithErrors(nil, foo, nil)

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "error",
			Errors: []errors.Assertion{
				&serrors.Assertion{
					Message: "foo",
				},
			},
		}, err)
	})
}
