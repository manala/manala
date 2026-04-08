package manifest_test

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/recipe/manifest"
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
	repositoryBaseURL := filepath.FromSlash("testdata/LoaderSuite/TestHandlerErrors")
	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler)),
	))

	tests := []struct {
		test     string
		expected errors.Assertion
	}{
		{
			test: "Directory",
			expected: &serrors.Assertion{
				Type:    serrors.Error{},
				Message: "recipe manifest is a directory",
				Arguments: []any{
					"dir", filepath.Join(repositoryBaseURL, "Directory", "repository", "recipe", ".manala.yaml"),
				},
			},
		},
		{
			test: "SyntaxError",
			expected: &serrors.Assertion{
				Type:    serrors.Error{},
				Message: "unable to parse recipe manifest",
				Arguments: []any{
					"file", filepath.Join(repositoryBaseURL, "SyntaxError", "repository", "recipe", ".manala.yaml"),
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
				Type:    serrors.Error{},
				Message: "unable to parse recipe manifest",
				Arguments: []any{
					"file", filepath.Join(repositoryBaseURL, "Empty", "repository", "recipe", ".manala.yaml"),
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
				Type:    serrors.Error{},
				Message: "unable to parse recipe manifest",
				Arguments: []any{
					"file", filepath.Join(repositoryBaseURL, "MultipleDocuments", "repository", "recipe", ".manala.yaml"),
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
				Type:    serrors.Error{},
				Message: "unable to parse recipe manifest",
				Arguments: []any{
					"file", filepath.Join(repositoryBaseURL, "NotMap", "repository", "recipe", ".manala.yaml"),
					"line", 1, "column", 1,
				},
				Details: `
					> 1 | foo
					      ^
					* yaml document must be a map
				`,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			repository, _ := repositoryLoader.Load(filepath.Join(repositoryBaseURL, test.test, "repository"))

			chainMock := &recipe.LoaderHandlerChainMock{}

			handler := manifest.NewLoaderHandler(slog.New(slog.DiscardHandler))
			recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

			s.Nil(recipe)
			errors.Equal(s.T(), test.expected, err)
			chainMock.AssertExpectations(s.T())
		})
	}
}

func (s *LoaderSuite) TestHandler() {
	repositoryURL := filepath.FromSlash("testdata/LoaderSuite/TestHandler/repository")

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler)),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	chainMock := &recipe.LoaderHandlerChainMock{}

	handler := manifest.NewLoaderHandler(slog.New(slog.DiscardHandler))
	recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

	s.Require().NoError(err)
	s.Equal(filepath.Join(repositoryURL, "recipe"), recipe.Dir())
	s.Equal(repositoryURL, recipe.Repository().URL())
	chainMock.AssertExpectations(s.T())
}
