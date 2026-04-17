package manifest_test

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/expect"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestHandlerErrors() {
	repositoryBaseURL := filepath.FromSlash("testdata/LoaderSuite/TestHandlerErrors")
	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.New(io.Discard)),
	))

	tests := []struct {
		test     string
		expected expect.ErrorExpectation
	}{
		{
			test: "Directory",
			expected: serrors.Expectation{
				Message: "recipe manifest is a directory",
				Attrs: [][2]any{
					{"dir", filepath.Join(repositoryBaseURL, "Directory", "repository", "recipe", ".manala.yaml")},
				},
			},
		},
		{
			test: "SyntaxError",
			expected: serrors.Expectation{
				Message: "unable to parse recipe manifest",
				Dump: heredoc.Doc(`
					in %[1]s:1:1
					> 1 | @
					      ^
					* '@' is a reserved character
				`,
					filepath.Join(repositoryBaseURL, "SyntaxError", "repository", "recipe", ".manala.yaml"),
				),
			},
		},
		{
			test: "Empty",
			expected: serrors.Expectation{
				Message: "unable to parse recipe manifest",
				Dump: heredoc.Doc(`
					in %[1]s:0
					* empty yaml content
				`,
					filepath.Join(repositoryBaseURL, "Empty", "repository", "recipe", ".manala.yaml"),
				),
			},
		},
		{
			test: "MultipleDocuments",
			expected: serrors.Expectation{
				Message: "unable to parse recipe manifest",
				Dump: heredoc.Doc(`
					in %[1]s:5:1
					  2 | 
					  3 | document: 1
					  4 | 
					> 5 | ---
					      ^
					  6 | 
					  7 | document: 2
					* multiple documents yaml content
				`,
					filepath.Join(repositoryBaseURL, "MultipleDocuments", "repository", "recipe", ".manala.yaml"),
				),
			},
		},
		{
			test: "NotMap",
			expected: serrors.Expectation{
				Message: "unable to parse recipe manifest",
				Dump: heredoc.Doc(`
					in %[1]s:1:1
					> 1 | foo
					      ^
					* yaml document must be a map
				`,
					filepath.Join(repositoryBaseURL, "NotMap", "repository", "recipe", ".manala.yaml"),
				),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			repository, _ := repositoryLoader.Load(filepath.Join(repositoryBaseURL, test.test, "repository"))

			chainMock := &recipe.LoaderHandlerChainMock{}

			handler := manifest.NewLoaderHandler(log.New(io.Discard))
			recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

			s.Nil(recipe)
			expect.Error(s.T(), test.expected, err)
			chainMock.AssertExpectations(s.T())
		})
	}
}

func (s *LoaderSuite) TestHandler() {
	repositoryURL := filepath.FromSlash("testdata/LoaderSuite/TestHandler/repository")

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.New(io.Discard)),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	chainMock := &recipe.LoaderHandlerChainMock{}

	handler := manifest.NewLoaderHandler(log.New(io.Discard))
	recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

	s.Require().NoError(err)
	s.Equal(filepath.Join(repositoryURL, "recipe"), recipe.Dir())
	s.Equal(repositoryURL, recipe.Repository().URL())
	chainMock.AssertExpectations(s.T())
}
