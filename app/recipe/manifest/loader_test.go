package manifest_test

import (
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

func (s *LoaderSuite) TestHandlerErrors() {
	dir := filepath.FromSlash("testdata/LoaderSuite/TestHandlerErrors")
	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
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
					{"dir", filepath.Join(dir, "Directory", "repository", "recipe", ".manala.yaml")},
				},
			},
		},
		{
			test: "UnparsableAnnotation",
			expected: serrors.Expectation{
				Message: "unable to parse recipe manifest",
				Dump: heredoc.Doc(`
				at %[1]s:4:12

				  1 │ manala:
				  2 │   description: description
				  3 │
				▶ 4 │ # @schema foo
				    ├────────────╯ invalid character 'o' in literal false (expecting 'a')
				  5 │ node: ~
			`,
					filepath.Join(dir, "UnparsableAnnotation", "repository", "recipe", ".manala.yaml"),
				),
			},
		},
		{
			test: "InvalidAnnotation",
			expected: serrors.Expectation{
				Message: "unable to parse recipe manifest",
				Dump: heredoc.Doc(`
				at %[1]s:4:11

				  1 │ manala:
				  2 │   description: description
				  3 │
				▶ 4 │ # @option {
				    ├───────────╯ missing option label property
				  5 │ #   "foo": "bar"
				  6 │ # }
				  7 │ node: foo
			`,
					filepath.Join(dir, "InvalidAnnotation", "repository", "recipe", ".manala.yaml"),
				),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			repository, _ := repositoryLoader.Load(filepath.Join(dir, test.test, "repository"))

			chainMock := &recipe.LoaderHandlerChainMock{}

			handler := manifest.NewLoaderHandler(log.Discard)
			recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)

			s.Nil(recipe)
			expect.Error(s.T(), test.expected, err)
			chainMock.AssertExpectations(s.T())
		})
	}
}
