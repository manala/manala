package serrors

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ArgumentsSuite struct{ suite.Suite }

func TestArgumentsSuite(t *testing.T) {
	suite.Run(t, new(ArgumentsSuite))
}

func (s *ArgumentsSuite) Test() {
	args := NewArguments()

	s.Equal([]any(nil), args.ErrorArguments())

	args.AppendArguments("foo", "bar")

	s.Equal([]any{"foo", "bar"}, args.ErrorArguments())

	args.PrependArguments("baz", "qux")

	s.Equal([]any{"baz", "qux", "foo", "bar"}, args.ErrorArguments())
}
