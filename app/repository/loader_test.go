package repository

import (
	"manala/app"
	"manala/internal/serrors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestLoadErrors() {
	loader := NewLoader()

	s.Run("NotFound", func() {
		repository, err := loader.Load("url")

		s.Nil(repository)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", "url",
			},
		}, err)
	})
}

func (s *LoaderSuite) TestLoad() {
	repositoryMock := &app.RepositoryMock{}

	handlerMock := &LoaderHandlerMock{}
	handlerMock.
		On("Handle", &LoaderQuery{URL: "url"}, mock.Anything).Return(repositoryMock, nil)

	loader := NewLoader(WithLoaderHandlers(handlerMock))

	repository, err := loader.Load("url")

	s.Require().NoError(err)
	s.Equal(repositoryMock, repository)
	handlerMock.AssertExpectations(s.T())
}
