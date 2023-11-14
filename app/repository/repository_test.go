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
	repository := NewRepository(
		"url",
		"dir",
	)

	s.Equal("url", repository.Url())
	s.Equal("dir", repository.Dir())
}
