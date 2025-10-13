package manifest_test

import (
	"path/filepath"
	"testing"

	"manala/app/recipe"
	"manala/app/recipe/manifest"
	"manala/app/repository"
	"manala/app/repository/getter"
	"manala/internal/log"
	"manala/internal/serrors"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestHandlerErrors() {
	s.Run("Directory", func() {
		repositoryURL := filepath.FromSlash("testdata/LoaderSuite/TestHandlerErrors/Directory/repository")

		repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
			getter.NewFileLoaderHandler(log.Discard),
		))
		repository, _ := repositoryLoader.Load(repositoryURL)

		chainMock := &recipe.LoaderHandlerChainMock{}

		handler := manifest.NewLoaderHandler(log.Discard)
		recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

		s.Nil(recipe)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(repositoryURL, "recipe", ".manala.yaml"),
			},
		}, err)
		chainMock.AssertExpectations(s.T())
	})
}

func (s *LoaderSuite) TestHandler() {
	repositoryURL := filepath.FromSlash("testdata/LoaderSuite/TestHandler/repository")

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	chainMock := &recipe.LoaderHandlerChainMock{}

	handler := manifest.NewLoaderHandler(log.Discard)
	recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

	s.Require().NoError(err)
	s.Equal(filepath.Join(repositoryURL, "recipe"), recipe.Dir())
	s.Equal(repositoryURL, recipe.Repository().URL())
	chainMock.AssertExpectations(s.T())
}
