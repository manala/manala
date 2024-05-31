package accessor

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGet() {
	var value any = "foo"

	accessor := New(&value)
	_value, err := accessor.Get()

	s.Require().NoError(err)
	s.Equal("foo", _value)
}

func (s *Suite) TestSet() {
	var value any = "foo"

	accessor := New(&value)
	err := accessor.Set("bar")

	s.Require().NoError(err)
	s.Equal("bar", value)
}
