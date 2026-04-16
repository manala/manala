package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/project"
	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/app/recipe"
	recipeManifest "github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/expect"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestHandler() {
	projectDir := filepath.FromSlash("testdata/LoaderSuite/TestHandler/project")

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
	))
	recipeLoader := recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(
		recipeManifest.NewLoaderHandler(log.Discard),
	))

	chainMock := &project.LoaderHandlerChainMock{}

	handler := manifest.NewLoaderHandler(log.Discard, repositoryLoader, recipeLoader)
	project, err := handler.Handle(&project.LoaderQuery{Dir: projectDir}, chainMock)

	s.Require().NoError(err)
	s.Equal(projectDir, project.Dir())
	s.Equal(map[string]any{"foo": "baz"}, project.Vars())
	chainMock.AssertExpectations(s.T())
}

func (s *LoaderSuite) TestHandlerErrors() {
	dir := filepath.FromSlash("testdata/LoaderSuite/TestHandlerErrors")
	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
	))
	recipeLoader := recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(
		recipeManifest.NewLoaderHandler(log.Discard),
	))

	tests := []struct {
		test     string
		expected expect.ErrorExpectation
	}{
		{
			test: "Directory",
			expected: serrors.Expectation{
				Message: "project manifest is a directory",
				Attrs: [][2]any{
					{"dir", filepath.Join(dir, "Directory", "project", ".manala.yaml")},
				},
			},
		},
		{
			test: "InvalidVars",
			expected: serrors.Expectation{
				Message: "invalid project manifest vars",
				Attrs: [][2]any{
					{"file", filepath.Join(dir, "InvalidVars", "project", ".manala.yaml")},
				},
				Errors: []expect.ErrorExpectation{
					serrors.Expectation{
						Message: "invalid type",
						Attrs: [][2]any{
							{"expected", "integer"},
							{"actual", "string"},
							{"path", "foo"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			chainMock := &project.LoaderHandlerChainMock{}

			handler := manifest.NewLoaderHandler(log.Discard, repositoryLoader, recipeLoader)
			project, err := handler.Handle(&project.LoaderQuery{Dir: filepath.Join(dir, test.test, "project")}, chainMock)

			s.Nil(project)
			expect.Error(s.T(), test.expected, err)
			chainMock.AssertExpectations(s.T())
		})
	}
}
