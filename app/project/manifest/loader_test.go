package manifest

import (
	"github.com/stretchr/testify/suite"
	"manala/app/project"
	"manala/app/recipe"
	"manala/app/recipe/manifest"
	"manala/app/repository"
	"manala/app/repository/getter"
	"manala/internal/filepath/filter"
	"manala/internal/log"
	"manala/internal/serrors"
	"path/filepath"
	"testing"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestLoaderHandlerErrors() {
	s.Run("Directory", func() {
		projectDir := filepath.FromSlash("testdata/LoaderSuite/TestLoaderHandlerErrors/Directory/project")

		repositoryLoader := repository.NewLoader(log.Discard)
		recipeLoader := recipe.NewLoader(log.Discard, filter.New())

		chainMock := &project.LoaderHandlerChainMock{}

		handler := NewLoaderHandler(log.Discard, repositoryLoader, recipeLoader)
		project, err := handler.Handle(&project.LoaderQuery{Dir: projectDir}, chainMock)

		s.Nil(project)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projectDir, ".manala.yaml"),
			},
		}, err)
		chainMock.AssertExpectations(s.T())
	})
	s.Run("Vars", func() {
		projectDir := filepath.FromSlash("testdata/LoaderSuite/TestLoaderHandlerErrors/Vars/project")

		repositoryLoader := repository.NewLoader(log.Discard, getter.NewFileLoaderHandler(log.Discard))
		recipeLoader := recipe.NewLoader(log.Discard, filter.New(), manifest.NewLoaderHandler(log.Discard))

		chainMock := &project.LoaderHandlerChainMock{}

		handler := NewLoaderHandler(log.Discard, repositoryLoader, recipeLoader)
		project, err := handler.Handle(&project.LoaderQuery{Dir: projectDir}, chainMock)

		s.Nil(project)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
			Message: "invalid project manifest vars",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
			},
			Errors: []*serrors.Assertion{
				{
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
						   3 |   repository: testdata/LoaderSuite/TestLoaderHandlerErrors/Vars/repository
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

func (s *LoaderSuite) TestLoaderHandler() {
	projectDir := filepath.FromSlash("testdata/LoaderSuite/TestLoaderHandler/project")

	repositoryLoader := repository.NewLoader(log.Discard, getter.NewFileLoaderHandler(log.Discard))
	recipeLoader := recipe.NewLoader(log.Discard, filter.New(), manifest.NewLoaderHandler(log.Discard))

	chainMock := &project.LoaderHandlerChainMock{}

	handler := NewLoaderHandler(log.Discard, repositoryLoader, recipeLoader)
	project, err := handler.Handle(&project.LoaderQuery{Dir: projectDir}, chainMock)

	s.Equal(projectDir, project.Dir())
	s.Equal(map[string]any{"foo": "bar"}, project.Vars())
	s.NoError(err)
	chainMock.AssertExpectations(s.T())
}
