package manifest_test

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/project"
	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/app/recipe"
	recipeManifest "github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestHandlerErrors() {
	s.Run("Directory", func() {
		projectDir := filepath.FromSlash("testdata/LoaderSuite/TestHandlerErrors/Directory/project")

		repositoryLoader := repository.NewLoader()
		recipeLoader := recipe.NewLoader(slog.New(slog.DiscardHandler))

		chainMock := &project.LoaderHandlerChainMock{}

		handler := manifest.NewLoaderHandler(slog.New(slog.DiscardHandler), repositoryLoader, recipeLoader)
		project, err := handler.Handle(&project.LoaderQuery{Dir: projectDir}, chainMock)

		s.Nil(project)
		errors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projectDir, ".manala.yaml"),
			},
		}, err)
		chainMock.AssertExpectations(s.T())
	})
	s.Run("Vars", func() {
		projectDir := filepath.FromSlash("testdata/LoaderSuite/TestHandlerErrors/Vars/project")

		repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
			getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler)),
		))
		recipeLoader := recipe.NewLoader(slog.New(slog.DiscardHandler), recipe.WithLoaderHandlers(
			recipeManifest.NewLoaderHandler(slog.New(slog.DiscardHandler)),
		))

		chainMock := &project.LoaderHandlerChainMock{}

		handler := manifest.NewLoaderHandler(slog.New(slog.DiscardHandler), repositoryLoader, recipeLoader)
		project, err := handler.Handle(&project.LoaderQuery{Dir: projectDir}, chainMock)

		s.Nil(project)
		errors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "invalid project manifest vars",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
			},
			Errors: []errors.Assertion{
				&serrors.Assertion{
					Type:    serrors.Error{},
					Message: "invalid type",
					Arguments: []any{
						"expected", "integer",
						"actual", "string",
						"path", "foo",
						"line", 5,
						"column", 6,
					},
					Details: `
						   2 |   recipe: recipe
						   3 |   repository: testdata/LoaderSuite/TestHandlerErrors/Vars/repository
						   4 |
						>  5 | foo: bar
						            ^
					`,
				},
			},
		}, err)
		chainMock.AssertExpectations(s.T())
	})
}

func (s *LoaderSuite) TestHandler() {
	projectDir := filepath.FromSlash("testdata/LoaderSuite/TestHandler/project")

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler)),
	))
	recipeLoader := recipe.NewLoader(slog.New(slog.DiscardHandler), recipe.WithLoaderHandlers(
		recipeManifest.NewLoaderHandler(slog.New(slog.DiscardHandler)),
	))

	chainMock := &project.LoaderHandlerChainMock{}

	handler := manifest.NewLoaderHandler(slog.New(slog.DiscardHandler), repositoryLoader, recipeLoader)
	project, err := handler.Handle(&project.LoaderQuery{Dir: projectDir}, chainMock)

	s.Require().NoError(err)
	s.Equal(projectDir, project.Dir())
	s.Equal(map[string]any{"foo": "bar"}, project.Vars())
	chainMock.AssertExpectations(s.T())
}
