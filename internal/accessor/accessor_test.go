package accessor_test

import (
	"testing"

	"github.com/manala/manala/internal/accessor"

	"github.com/stretchr/testify/suite"
)

type AccessorSuite struct{ suite.Suite }

func TestAccessorSuite(t *testing.T) {
	suite.Run(t, new(AccessorSuite))
}

func (s *AccessorSuite) TestGet() {
	var value any = "foo"

	accessor := accessor.New(&value)
	_value, err := accessor.Get()

	s.Require().NoError(err)
	s.Equal("foo", _value)
}

func (s *AccessorSuite) TestSet() {
	var value any = "foo"

	accessor := accessor.New(&value)
	err := accessor.Set("bar")

	s.Require().NoError(err)
	s.Equal("bar", value)
}
