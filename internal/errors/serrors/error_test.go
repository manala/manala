package serrors

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorSuite struct{ suite.Suite }

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (s *ErrorSuite) Test() {
	err := New("message")

	Equal(s.Assert(), &Assert{
		Type:    &Error{},
		Message: "message",
	}, err)

	err = err.WithArguments(
		"foo", "bar",
	)

	Equal(s.Assert(), &Assert{
		Type:    &Error{},
		Message: "message",
		Arguments: []any{
			"foo", "bar",
		},
	}, err)
}
