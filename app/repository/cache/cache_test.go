package cache

import (
	"github.com/stretchr/testify/suite"
	"manala/app"
	"testing"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test() {
	cache := New()

	repository, ok := cache.Get("foo")

	s.Nil(repository)
	s.False(ok)

	repositoryMock := &app.RepositoryMock{}

	cache.Set("foo", repositoryMock)

	repository, ok = cache.Get("foo")

	s.Equal(repositoryMock, repository)
	s.True(ok)
}
