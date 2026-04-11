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
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestHandlerErrors() {
	projectBaseDir := filepath.FromSlash("testdata/LoaderSuite/TestHandlerErrors")
	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler)),
	))
	recipeLoader := recipe.NewLoader(slog.New(slog.DiscardHandler), recipe.WithLoaderHandlers(
		recipeManifest.NewLoaderHandler(slog.New(slog.DiscardHandler)),
	))

	tests := []struct {
		test     string
		expected errors.Assertion
	}{
		{
			test: "Directory",
			expected: &serrors.Assertion{
				Message: "project manifest is a directory",
				Arguments: []any{
					"dir", filepath.Join(projectBaseDir, "Directory", "project", ".manala.yaml"),
				},
			},
		},
		{
			test: "SyntaxError",
			expected: &serrors.Assertion{
				Message: "unable to parse project manifest",
				Arguments: []any{
					"file", filepath.Join(projectBaseDir, "SyntaxError", "project", ".manala.yaml"),
					"line", 1, "column", 1,
				},
				Details: `
					> 1 | @
					      ^
					* '@' is a reserved character
				`,
			},
		},
		{
			test: "Empty",
			expected: &serrors.Assertion{
				Message: "unable to parse project manifest",
				Arguments: []any{
					"file", filepath.Join(projectBaseDir, "Empty", "project", ".manala.yaml"),
				},
				Errors: []errors.Assertion{
					&parsing.Assertion{
						Err: &serrors.Assertion{
							Message: "empty yaml content",
						},
					},
				},
			},
		},
		{
			test: "MultipleDocuments",
			expected: &serrors.Assertion{
				Message: "unable to parse project manifest",
				Arguments: []any{
					"file", filepath.Join(projectBaseDir, "MultipleDocuments", "project", ".manala.yaml"),
					"line", 5, "column", 1,
				},
				Details: `
					  3 | document: 1
					  4 |
					> 5 | ---
					      ^
					  6 |
					  7 | document: 2
					* multiple documents yaml content
				`,
			},
		},
		{
			test: "NotMap",
			expected: &serrors.Assertion{
				Message: "unable to parse project manifest",
				Arguments: []any{
					"file", filepath.Join(projectBaseDir, "NotMap", "project", ".manala.yaml"),
					"line", 1, "column", 1,
				},
				Details: `
					> 1 | foo
					      ^
					* yaml document must be a map
				`,
			},
		},
		{
			test: "Vars",
			expected: &serrors.Assertion{
				Message: "invalid project manifest vars",
				Arguments: []any{
					"file", filepath.Join(projectBaseDir, "Vars", "project", ".manala.yaml"),
				},
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "invalid type",
						Arguments: []any{
							"expected", "integer",
							"actual", "string",
							"path", "foo",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			chainMock := &project.LoaderHandlerChainMock{}

			handler := manifest.NewLoaderHandler(slog.New(slog.DiscardHandler), repositoryLoader, recipeLoader)
			project, err := handler.Handle(&project.LoaderQuery{Dir: filepath.Join(projectBaseDir, test.test, "project")}, chainMock)

			s.Nil(project)
			errors.Equal(s.T(), test.expected, err)
			chainMock.AssertExpectations(s.T())
		})
	}
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
	s.Equal(map[string]any{"foo": "baz"}, project.Vars())
	chainMock.AssertExpectations(s.T())
}
