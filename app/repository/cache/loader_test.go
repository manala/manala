package cache

import (
	"manala/app"
	"manala/app/repository"
	"manala/internal/log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestHandler() {
	cache := New()

	handler := NewLoaderHandler(log.Discard, cache)
	handlerQueryFoo := &repository.LoaderQuery{Url: "foo"}
	handlerQueryBar := &repository.LoaderQuery{Url: "bar"}

	repositoryMock := &app.RepositoryMock{}

	chainMock := &repository.LoaderHandlerChainMock{}

	// First call should chain to next handler
	chainMockCall := chainMock.
		On("Next", handlerQueryFoo).Return(repositoryMock, nil)

	repository, err := handler.Handle(handlerQueryFoo, chainMock)

	s.Equal(repositoryMock, repository)
	s.NoError(err)
	chainMock.AssertExpectations(s.T())

	// Second call with same name should extract from cache, and not chain to next handler
	chainMockCall.Unset()

	repository, err = handler.Handle(handlerQueryFoo, chainMock)

	s.Equal(repositoryMock, repository)
	s.NoError(err)
	chainMock.AssertExpectations(s.T())

	// Third call with different name should chain to next handler
	chainMock.
		On("Next", handlerQueryBar).Return(repositoryMock, nil)

	repository, err = handler.Handle(handlerQueryBar, chainMock)

	s.Equal(repositoryMock, repository)
	s.NoError(err)
	chainMock.AssertExpectations(s.T())
}
