package filepath

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExcludeSuite struct{ suite.Suite }

func TestExcludeSuite(t *testing.T) {
	suite.Run(t, new(ExcludeSuite))
}

func (s *ExcludeSuite) Test() {
	s.True(Exclude(".git"))
	s.True(Exclude("foo/.git"))
	s.False(Exclude("foo"))
}
