package repository

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"manala/app/mocks"
	internalLog "manala/internal/log"
	"testing"
)

type CacheManagerSuite struct{ suite.Suite }

func TestCacheManagerSuite(t *testing.T) {
	suite.Run(t, new(CacheManagerSuite))
}

func (s *CacheManagerSuite) TestLoadRepositoryErrors() {
	log := internalLog.New(io.Discard)

	cascadingRepoMock := mocks.MockRepository()
	cascadingError := errors.New("error")

	cascadingManagerMock := mocks.MockRepositoryManager()
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(cascadingRepoMock, cascadingError)

	manager := NewCacheManager(
		log,
		cascadingManagerMock,
	)
	repo, err := manager.LoadRepository("url")

	s.Nil(repo)
	s.ErrorIs(err, cascadingError)
}

func (s *CacheManagerSuite) TestLoadRepository() {
	log := internalLog.New(io.Discard)

	cascadingRepoMock := mocks.MockRepository()

	cascadingManagerMock := mocks.MockRepositoryManager()
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(mocks.MockRepository(), nil)

	manager := NewCacheManager(
		log,
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
