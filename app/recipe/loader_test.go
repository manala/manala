package recipe_test

import (
	"manala/app"
	"manala/app/recipe"
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
	loader := recipe.NewLoader(log.Discard)

	s.Run("NotFound", func() {
		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("URL").Return("url")

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

	handlerMock := &recipe.LoaderHandlerMock{}
	handlerMock.
		On("Handle", &recipe.LoaderQuery{Repository: repositoryMock, Name: "name"}, mock.Anything).Return(recipeMock, nil)

	loader := recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(handlerMock))

	recipe, err := loader.Load(repositoryMock, "name")

	s.Require().NoError(err)
	s.Equal(recipeMock, recipe)
	handlerMock.AssertExpectations(s.T())
}

func (s *LoaderSuite) TestLoadAllErrors() {
	loader := recipe.NewLoader(log.Discard)

	s.Run("EmptyRepository", func() {
		repositoryURL := filepath.FromSlash("testdata/LoaderSuite/TestLoadAllErrors/EmptyRepository/repository")

		repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
			getter.NewFileLoaderHandler(log.Discard),
		))
		repository, _ := repositoryLoader.Load(repositoryURL)

		recipes, err := loader.LoadAll(repository)

		s.Empty(recipes)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.EmptyRepositoryError{},
			Message: "empty repository",
			Arguments: []any{
				"url", repositoryURL,
			},
		}, err)
	})
}

func (s *LoaderSuite) TestLoadAll() {
	repositoryURL := filepath.FromSlash("testdata/LoaderSuite/TestLoadAll/repository")

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	recipeMock := &app.RecipeMock{}

	handlerMock := &recipe.LoaderHandlerMock{}
	handlerMock.
		On("Handle", &recipe.LoaderQuery{Repository: repository, Name: "foo"}, mock.Anything).Return(recipeMock, nil).
		On("Handle", &recipe.LoaderQuery{Repository: repository, Name: "bar"}, mock.Anything).Return(recipeMock, nil)

	loader := recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(handlerMock))

	recipes, err := loader.LoadAll(repository)

	s.Require().NoError(err)
	s.Equal([]app.Recipe{
		recipeMock,
		recipeMock,
	}, recipes)
	handlerMock.AssertExpectations(s.T())
}
