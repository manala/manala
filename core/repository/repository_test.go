package repository

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepositorySuite struct{ suite.Suite }

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) Test() {
	repo := NewRepository(
		"url",
		"dir",
	)

	s.Equal("url", repo.Url())
	s.Equal("dir", repo.Dir())
}
