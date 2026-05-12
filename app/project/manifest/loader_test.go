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
		expected expectation.ErrorExpectation
	}{
		{
			test: "Directory",
			expected: serror.Expectation{
				Msg: "project manifest is a directory",
				Attrs: [][2]any{
					{"dir", filepath.Join(dir, "Directory", "project", ".manala.yaml")},
				},
			},
		},
		{
			test: "Unparsable",
			expected: serror.Expectation{
				Msg: "unable to parse project manifest",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:1:1

						▶ 1 │ @
						    ├─╯ '@' is a reserved character
					`,
						filepath.Join(dir, "Unparsable", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "MissingConfig",
			expected: serror.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:1:4

						▶ 1 │ foo: bar
						    ├────╯ missing "manala" property
					`,
						filepath.Join(dir, "MissingConfig", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "UndecodableConfig",
			expected: serror.Expectation{
				Msg: "unable to decode project manifest config",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:1:9

						▶ 1 │ manala: foo
						    ├─────────╯ string was used where mapping is expected
					`,
						filepath.Join(dir, "UndecodableConfig", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "InvalidConfig",
			expected: serror.Expectation{
				Msg: "unable to decode project manifest config",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:1:9

						▶ 1 │ manala: {}
						    ├─────────╯ missing property 'recipe'
					`,
						filepath.Join(dir, "InvalidConfig", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "InvalidVars",
			expected: serror.Expectation{
				Msg: "invalid project manifest vars",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`
						at %[1]s:5:6

						  2 │   recipe: recipe
						  3 │   repository: testdata/LoaderSuite/TestHandlerErrors/InvalidVars/repository
						  4 │
						▶ 5 │ foo: bar
						    ├──────╯ got string, want integer
					`,
						filepath.Join(dir, "InvalidVars", "project", ".manala.yaml"),
					)),
				),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			chainMock := &project.LoaderHandlerChainMock{}

			handler := manifest.NewLoaderHandler(log.Discard, repositoryLoader, recipeLoader)
			project, err := handler.Handle(&project.LoaderQuery{Dir: filepath.Join(dir, test.test, "project")}, chainMock)

			s.Nil(project)
			expectation.ExpectError(s.T(), test.expected, err)
			chainMock.AssertExpectations(s.T())
		})
	}
}
