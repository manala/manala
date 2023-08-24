package serrors

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type WrapErrorSuite struct{ suite.Suite }

func TestWrapErrorSuite(t *testing.T) {
	suite.Run(t, new(WrapErrorSuite))
}

func (s *WrapErrorSuite) Test() {
	err := Wrap("message", New("wrap"))

	Equal(s.Assert(), &Assert{
		Type:    &WrapError{},
		Message: "message",
		Error: &Assert{
			Type:    &Error{},
			Message: "wrap",
		},
	}, err)

	err = err.WithArguments(
		"foo", "bar",
	)

	Equal(s.Assert(), &Assert{
		Type:    &WrapError{},
		Message: "message",
		Arguments: []any{
			"foo", "bar",
		},
		Error: &Assert{
			Type:    &Error{},
			Message: "wrap",
		},
	}, err)
}
