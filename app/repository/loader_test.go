package repository

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"manala/app"
	"manala/internal/serrors"
	"testing"
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
		On("Handle", &LoaderQuery{Url: "url"}, mock.Anything).Return(repositoryMock, nil)

	loader := NewLoader(WithLoaderHandlers(handlerMock))

	repository, err := loader.Load("url")

	s.Equal(repositoryMock, repository)
	s.NoError(err)
	handlerMock.AssertExpectations(s.T())
}
