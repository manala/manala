package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/source"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/testing/expectation"
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
		expected expectation.ErrorExpectation
	}{
		{
			test: "Directory",
			expected: serror.Expectation{
				Msg: "recipe manifest is a directory",
				Attrs: [][2]any{
					{"dir", filepath.Join(dir, "Directory", "repository", "recipe", ".manala.yaml")},
				},
			},
		},
		{
			test: "Unparsable",
			expected: serror.Expectation{
				Msg: "unable to parse recipe manifest",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:1:1

						▶ 1 │ @
						    ├─╯ '@' is a reserved character
					`,
						filepath.Join(dir, "Unparsable", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "MissingConfig",
			expected: serror.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:1:4

						▶ 1 │ foo: bar
						    ├────╯ missing "manala" property
					`,
						filepath.Join(dir, "MissingConfig", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "UndecodableConfig",
			expected: serror.Expectation{
				Msg: "unable to decode recipe manifest config",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:1:9

						▶ 1 │ manala: foo
						    ├─────────╯ string was used where mapping is expected
					`,
						filepath.Join(dir, "UndecodableConfig", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "InvalidConfig",
			expected: serror.Expectation{
				Msg: "unable to decode recipe manifest config",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:1:9

						▶ 1 │ manala: {}
						    ├─────────╯ missing property 'description'
					`,
						filepath.Join(dir, "InvalidConfig", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "UnparsableAnnotation",
			expected: serror.Expectation{
				Msg: "unable to infer recipe manifest vars",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:4:12

						  1 │ manala:
						  2 │   description: description
						  3 │
						▶ 4 │ # @schema foo
						    ├────────────╯ invalid character 'o' in literal false (expecting 'a')
						  5 │ node: ~
					`,
						filepath.Join(dir, "UnparsableAnnotation", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "InvalidAnnotation",
			expected: serror.Expectation{
				Msg: "unable to infer recipe manifest vars",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:4:11

						  1 │ manala:
						  2 │   description: description
						  3 │
						▶ 4 │ # @option {
						    ├───────────╯ missing property 'label'
						  5 │ #   "foo": "bar"
						  6 │ # }
						  7 │ node: foo
					`,
						filepath.Join(dir, "InvalidAnnotation", "repository", "recipe", ".manala.yaml"),
					)),
					source.Expectation(heredoc.Doc(`
						at %[1]s:4:11

						  1 │ manala:
						  2 │   description: description
						  3 │
						▶ 4 │ # @option {
						    ├───────────╯ additional properties 'foo' not allowed
						  5 │ #   "foo": "bar"
						  6 │ # }
						  7 │ node: foo
					`,
						filepath.Join(dir, "InvalidAnnotation", "repository", "recipe", ".manala.yaml"),
					)),
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
			expectation.ExpectError(s.T(), test.expected, err)
			chainMock.AssertExpectations(s.T())
		})
	}
}
