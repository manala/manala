package getter_test

import (
	"testing"

	"manala/app/repository/getter"

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
