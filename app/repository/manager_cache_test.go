package repository

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app"
	"manala/internal/serrors"
	"testing"
)

type CacheManagerSuite struct{ suite.Suite }

func TestCacheManagerSuite(t *testing.T) {
	suite.Run(t, new(CacheManagerSuite))
}

func (s *CacheManagerSuite) TestLoadRepositoryErrors() {
	cascadingMock := &app.RepositoryMock{}
	cascadingError := serrors.New("error")

	cascadingManagerMock := &app.RepositoryManagerMock{}
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(cascadingMock, cascadingError)

	manager := NewCacheManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cascadingManagerMock,
	)
	repository, err := manager.LoadRepository("url")

	s.Nil(repository)
	s.Equal(err, cascadingError)
}

func (s *CacheManagerSuite) TestLoadRepository() {
	cascadingMock := &app.RepositoryMock{}

	cascadingManagerMock := &app.RepositoryManagerMock{}
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(&app.RepositoryMock{}, nil)

	manager := NewCacheManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cascadingManagerMock,
	)

	// First call should pass to embedded manager
	repository, err := manager.LoadRepository("foo")

	s.NoError(err)
	s.Equal(cascadingMock, repository)

	cascadingManagerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 1)

	// Second call with same name should extract from cache, and not pass to embedded manager
	repository, err = manager.LoadRepository("foo")

	s.NoError(err)
	s.Equal(cascadingMock, repository)

	cascadingManagerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 1)

	// Third call with different name should pass to embedded manager
	repository, err = manager.LoadRepository("bar")

	s.NoError(err)
	s.Equal(cascadingMock, repository)

	cascadingManagerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 2)
}
