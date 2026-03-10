package serrors_test

import (
	"fmt"
	"github.com/manala/manala/internal/serrors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestError() {
	s.Run("New", func() {
		err := serrors.New("error")

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "error",
		}, err)
	})

	s.Run("Message", func() {
		message := "message"

		err := serrors.New("error").
			WithMessage(message)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: message,
		}, err)
	})

	s.Run("Arguments", func() {
		foo := "foo"
		bar := "bar"

		err := serrors.New("error").
			WithArguments(foo, bar)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "error",
			Arguments: []any{
				foo, bar,
			},
		}, err)
	})

	s.Run("Details", func() {
		details := "details"

		err := serrors.New("error").
			WithDetails(details)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "error",
			Details: details,
		}, err)
	})

	s.Run("DetailsFunc", func() {
		detailsFunc := func(ansi bool) string {
			return fmt.Sprintf("details func %v", ansi)
		}

		err := serrors.New("error").
			WithDetailsFunc(detailsFunc)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "error",
			Details: "details func false",
		}, err)
	})

	s.Run("Errors", func() {
		foo := serrors.New("foo")
		bar := serrors.New("bar")

		err := serrors.New("error").
			WithErrors(foo, bar)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "error",
			Errors: []*serrors.Assertion{
				{
					Type:    serrors.Error{},
					Message: "foo",
				},
				{
					Type:    serrors.Error{},
					Message: "bar",
				},
			},
		}, err)
	})

	s.Run("NilErrors", func() {
		foo := serrors.New("foo")

		err := serrors.New("error").
			WithErrors(nil, foo, nil)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "error",
			Errors: []*serrors.Assertion{
				{
					Type:    serrors.Error{},
					Message: "foo",
				},
			},
		}, err)
	})
}
