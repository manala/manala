package repository

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app/mocks"
	"testing"
)

type CacheManagerSuite struct{ suite.Suite }

func TestCacheManagerSuite(t *testing.T) {
	suite.Run(t, new(CacheManagerSuite))
}

func (s *CacheManagerSuite) TestLoadRepositoryErrors() {
	cascadingRepoMock := &mocks.RepositoryMock{}
	cascadingError := errors.New("error")

	cascadingManagerMock := &mocks.RepositoryManagerMock{}
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(cascadingRepoMock, cascadingError)

	manager := NewCacheManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cascadingManagerMock,
	)
	repo, err := manager.LoadRepository("url")

	s.Nil(repo)

	s.ErrorIs(err, cascadingError)
}

func (s *CacheManagerSuite) TestLoadRepository() {
	cascadingRepoMock := &mocks.RepositoryMock{}

	cascadingManagerMock := &mocks.RepositoryManagerMock{}
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(&mocks.RepositoryMock{}, nil)

	manager := NewCacheManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cascadingManagerMock,
	)

	// First call should pass to embedded manager
	repo, err := manager.LoadRepository("foo")

	s.NoError(err)
	s.Equal(cascadingRepoMock, repo)

	cascadingManagerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 1)

	// Second call with same name should extract from cache, and not pass to embedded manager
	repo, err = manager.LoadRepository("foo")

	s.NoError(err)
	s.Equal(cascadingRepoMock, repo)

	cascadingManagerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 1)

	// Third call with different name should pass to embedded manager
	repo, err = manager.LoadRepository("bar")

	s.NoError(err)
	s.Equal(cascadingRepoMock, repo)

	cascadingManagerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 2)
}
