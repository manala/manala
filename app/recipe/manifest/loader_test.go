package manifest

import (
	"manala/app/recipe"
	"manala/app/repository"
	"manala/app/repository/getter"
	"manala/internal/log"
	"manala/internal/serrors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestLoaderHandlerErrors() {
	s.Run("Directory", func() {
		repositoryUrl := filepath.FromSlash("testdata/LoaderSuite/TestLoaderHandlerErrors/Directory/repository")

		repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
			getter.NewFileLoaderHandler(log.Discard),
		))
		repository, _ := repositoryLoader.Load(repositoryUrl)

		chainMock := &recipe.LoaderHandlerChainMock{}

		handler := NewLoaderHandler(log.Discard)
		recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

		s.Nil(recipe)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(repositoryUrl, "recipe", ".manala.yaml"),
			},
		}, err)
		chainMock.AssertExpectations(s.T())
	})
}

func (s *LoaderSuite) TestLoaderHandler() {
	repositoryUrl := filepath.FromSlash("testdata/LoaderSuite/TestLoaderHandler/repository")

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
	))
	repository, _ := repositoryLoader.Load(repositoryUrl)

	chainMock := &recipe.LoaderHandlerChainMock{}

	handler := NewLoaderHandler(log.Discard)
	recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

	s.Equal(filepath.Join(repositoryUrl, "recipe"), recipe.Dir())
	s.Equal(repositoryUrl, recipe.Repository().Url())
	s.NoError(err)
	chainMock.AssertExpectations(s.T())
}
