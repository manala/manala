package serrors

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type WrapsErrorSuite struct{ suite.Suite }

func TestWrapsErrorSuite(t *testing.T) {
	suite.Run(t, new(WrapsErrorSuite))
}

func (s *WrapsErrorSuite) Test() {
	err := Wraps("message",
		New("wrap 1"),
		New("wrap 2"),
	)

	Equal(s.Assert(), &Assert{
		Type:    &WrapsError{},
		Message: "message",
		Errors: []*Assert{
			{
				Type:    &Error{},
				Message: "wrap 1",
			},
			{
				Type:    &Error{},
				Message: "wrap 2",
			},
		},
	}, err)

	err = err.WithArguments(
		"foo", "bar",
	)

	Equal(s.Assert(), &Assert{
		Type:    &WrapsError{},
		Message: "message",
		Arguments: []any{
			"foo", "bar",
		},
		Errors: []*Assert{
			{
				Type:    &Error{},
				Message: "wrap 1",
			},
			{
				Type:    &Error{},
				Message: "wrap 2",
			},
		},
	}, err)
}
