package serrors

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestError() {
	s.Run("New", func() {
		err := New("error")

		Equal(s.Assert(), &Assert{
			Type:    Error{},
			Message: "error",
		}, err)
	})

	s.Run("Message", func() {
		message := "message"

		err := New("error").
			WithMessage(message)

		Equal(s.Assert(), &Assert{
			Type:    Error{},
			Message: message,
		}, err)
	})

	s.Run("Arguments", func() {
		foo := "foo"
		bar := "bar"

		err := New("error").
			WithArguments(foo, bar)

		Equal(s.Assert(), &Assert{
			Type:    Error{},
			Message: "error",
			Arguments: []any{
				foo, bar,
			},
		}, err)
	})

	s.Run("Details", func() {
		details := "details"

		err := New("error").
			WithDetails(details)

		Equal(s.Assert(), &Assert{
			Type:    Error{},
			Message: "error",
			Details: details,
		}, err)
	})

	s.Run("DetailsFunc", func() {
		detailsFunc := func(ansi bool) string {
			return fmt.Sprintf("details func %v", ansi)
		}

		err := New("error").
			WithDetailsFunc(detailsFunc)

		Equal(s.Assert(), &Assert{
			Type:    Error{},
			Message: "error",
			Details: "details func false",
		}, err)
	})

	s.Run("Errors", func() {
		foo := New("foo")
		bar := New("bar")

		err := New("error").
			WithErrors(foo, bar)

		Equal(s.Assert(), &Assert{
			Type:    Error{},
			Message: "error",
			Errors: []*Assert{
				{
					Type:    Error{},
					Message: "foo",
				},
				{
					Type:    Error{},
					Message: "bar",
				},
			},
		}, err)
	})
}
