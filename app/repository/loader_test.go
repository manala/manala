package repository_test

import (
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/testing/errors"
	"github.com/manala/manala/app/testing/mocks"
	"github.com/manala/manala/internal/testing/expect"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestLoadErrors() {
	loader := repository.NewLoader()

	s.Run("NotFound", func() {
		repository, err := loader.Load("url")

		s.Nil(repository)
		expect.Error(s.T(), errors.Expectation{
			Type:  &app.NotFoundRepositoryError{},
			Attrs: [][2]any{{"url", "url"}},
		}, err)
	})
}

func (s *LoaderSuite) TestLoad() {
	repositoryMock := &mocks.RepositoryMock{}

	handlerMock := &repository.LoaderHandlerMock{}
	handlerMock.
		On("Handle", &repository.LoaderQuery{URL: "url"}, mock.Anything).Return(repositoryMock, nil)

	loader := repository.NewLoader(repository.WithLoaderHandlers(handlerMock))

	repository, err := loader.Load("url")

	s.Require().NoError(err)
	s.Equal(repositoryMock, repository)
	handlerMock.AssertExpectations(s.T())
}
