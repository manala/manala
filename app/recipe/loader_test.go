package recipe

import (
	"manala/app"
	"manala/app/repository"
	"manala/app/repository/getter"
	"manala/internal/log"
	"manala/internal/serrors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestLoadErrors() {
	loader := NewLoader(log.Discard)

	s.Run("NotFound", func() {
		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("Url").Return("url")

		recipe, err := loader.Load(repositoryMock, "name")

		s.Nil(recipe)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRecipeError{},
			Message: "recipe not found",
			Arguments: []any{
				"repository", "url",
				"name", "name",
			},
		}, err)
	})
}

func (s *LoaderSuite) TestLoad() {
	repositoryMock := &app.RepositoryMock{}
	recipeMock := &app.RecipeMock{}

	handlerMock := &LoaderHandlerMock{}
	handlerMock.
		On("Handle", &LoaderQuery{Repository: repositoryMock, Name: "name"}, mock.Anything).Return(recipeMock, nil)

	loader := NewLoader(log.Discard, WithLoaderHandlers(handlerMock))

	recipe, err := loader.Load(repositoryMock, "name")

	s.Equal(recipeMock, recipe)
	s.NoError(err)
	handlerMock.AssertExpectations(s.T())
}

func (s *LoaderSuite) TestLoadAllErrors() {
	loader := NewLoader(log.Discard)

	s.Run("EmptyRepository", func() {
		repositoryUrl := filepath.FromSlash("testdata/LoaderSuite/TestLoadAllErrors/EmptyRepository/repository")

		repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
			getter.NewFileLoaderHandler(log.Discard),
		))
		repository, _ := repositoryLoader.Load(repositoryUrl)

		recipes, err := loader.LoadAll(repository)

		s.Empty(recipes)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.EmptyRepositoryError{},
			Message: "empty repository",
			Arguments: []any{
				"url", repositoryUrl,
			},
		}, err)
	})
}

func (s *LoaderSuite) TestLoadAll() {
	repositoryUrl := filepath.FromSlash("testdata/LoaderSuite/TestLoadAll/repository")

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
	))
	repository, _ := repositoryLoader.Load(repositoryUrl)

	recipeMock := &app.RecipeMock{}

	handlerMock := &LoaderHandlerMock{}
	handlerMock.
		On("Handle", &LoaderQuery{Repository: repository, Name: "foo"}, mock.Anything).Return(recipeMock, nil).
		On("Handle", &LoaderQuery{Repository: repository, Name: "bar"}, mock.Anything).Return(recipeMock, nil)

	loader := NewLoader(log.Discard, WithLoaderHandlers(handlerMock))

	recipes, err := loader.LoadAll(repository)

	s.Equal([]app.Recipe{
		recipeMock,
		recipeMock,
	}, recipes)
	s.NoError(err)
	handlerMock.AssertExpectations(s.T())
}
