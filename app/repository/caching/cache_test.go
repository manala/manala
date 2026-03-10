package caching_test

import (
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository/caching"

	"github.com/stretchr/testify/suite"
)

type CacheSuite struct{ suite.Suite }

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}

func (s *CacheSuite) Test() {
	cache := caching.NewCache()

	repository, ok := cache.Get("foo")

	s.Nil(repository)
	s.False(ok)

	repositoryMock := &app.RepositoryMock{}

	cache.Set("foo", repositoryMock)

	repository, ok = cache.Get("foo")

	s.Equal(repositoryMock, repository)
	s.True(ok)
}
