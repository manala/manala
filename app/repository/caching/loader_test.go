package caching_test

import (
	"io"
	"testing"

	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/caching"
	"github.com/manala/manala/app/testing/mocks"
	"github.com/manala/manala/internal/log"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestHandler() {
	cache := caching.NewCache()

	handler := caching.NewLoaderHandler(log.New(io.Discard), cache)
	handlerQueryFoo := &repository.LoaderQuery{URL: "foo"}
	handlerQueryBar := &repository.LoaderQuery{URL: "bar"}

	repositoryMock := &mocks.RepositoryMock{}

	chainMock := &repository.LoaderHandlerChainMock{}

	// First call should chain to next handler
	chainMockCall := chainMock.
		On("Next", handlerQueryFoo).Return(repositoryMock, nil)

	repository, err := handler.Handle(handlerQueryFoo, chainMock)

	s.Require().NoError(err)
	s.Equal(repositoryMock, repository)
	chainMock.AssertExpectations(s.T())

	// Second call with same name should extract from cache, and not chain to next handler
	chainMockCall.Unset()

	repository, err = handler.Handle(handlerQueryFoo, chainMock)

	s.Require().NoError(err)
	s.Equal(repositoryMock, repository)
	chainMock.AssertExpectations(s.T())

	// Third call with different name should chain to next handler
	chainMock.
		On("Next", handlerQueryBar).Return(repositoryMock, nil)

	repository, err = handler.Handle(handlerQueryBar, chainMock)

	s.Require().NoError(err)
	s.Equal(repositoryMock, repository)
	chainMock.AssertExpectations(s.T())
}
