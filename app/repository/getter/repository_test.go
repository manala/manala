package getter_test

import (
	"manala/app/repository/getter"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct{ suite.Suite }

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) Test() {
	repository := getter.NewRepository("url", "dir")

	s.Equal("url", repository.URL())
	s.Equal("dir", repository.Dir())
}
