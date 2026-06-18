package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/project"
	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/app/recipe"
	recipeManifest "github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/errors/serror/serrortest"
	"github.com/manala/manala/internal/errors/source/sourcetest"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestHandle() {
	projectDir := filepath.FromSlash("testdata/LoaderSuite/TestHandle/project")

	project, err := s.handle(projectDir)

	s.Require().NoError(err)

	s.Equal(projectDir, project.Dir())
	s.Equal(map[string]any{"foo": "baz"}, project.Vars())
}

func (s *LoaderSuite) TestHandleErrors() {
	dir := filepath.FromSlash("testdata/LoaderSuite/TestHandleErrors")

	tests := []struct {
		test     string
		expected expectation.ErrorExpectation
	}{
		{
			test: "Directory",
			expected: serrortest.Expectation{
				Msg: "project manifest is a directory",
				Attrs: [][2]any{
					{"dir", filepath.Join(dir, "Directory", "project", ".manala.yaml")},
				},
			},
		},
		{
			test: "Unparsable",
			expected: serrortest.Expectation{
				Msg: "unable to parse project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

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
			test: "Empty",
			expected: serrortest.Expectation{
				Msg: "unable to parse project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s

						  1 │

						empty yaml content
					`,
						filepath.Join(dir, "Empty", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "MultipleDocuments",
			expected: serrortest.Expectation{
				Msg: "unable to parse project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:2:1

						  1 │ foo: bar
						▶ 2 │ ---
						    ├─╯ multiple documents yaml content
						  3 │ foo: bar
					`,
						filepath.Join(dir, "MultipleDocuments", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "NotMap",
			expected: serrortest.Expectation{
				Msg: "unable to parse project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:1:1

						▶ 1 │ - foo
						    ├─╯ yaml document must be a map
						  2 │ - bar
					`,
						filepath.Join(dir, "NotMap", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "MapEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s

						  1 │ {}

						missing property 'manala'
					`,
						filepath.Join(dir, "MapEmpty", "project", ".manala.yaml"),
					)),
				),
			},
		},
		// Config
		{
			test: "ConfigMissing",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s

						  1 │ foo: bar

						missing property 'manala'
					`,
						filepath.Join(dir, "ConfigMissing", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigNotMap",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:1:9

						▶ 1 │ manala: foo
						    ├─────────╯ got string, want object
					`,
						filepath.Join(dir, "ConfigNotMap", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigMapEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:1:1

						▶ 1 │ manala: {}
						    ├─╯ missing property 'recipe'
					`,
						filepath.Join(dir, "ConfigMapEmpty", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigAdditionalProperty",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:3

						  1 │ manala:
						  2 │   recipe: recipe
						▶ 3 │   foo: bar
						    ├───╯ additional property 'foo' not allowed
					`,
						filepath.Join(dir, "ConfigAdditionalProperty", "project", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Recipe
		{
			test: "ConfigRecipeAbsent",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:1:1

						▶ 1 │ manala:
						    ├─╯ missing property 'recipe'
						  2 │   repository: repository
					`,
						filepath.Join(dir, "ConfigRecipeAbsent", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigRecipeNotString",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:2:11

						  1 │ manala:
						▶ 2 │   recipe: []
						    ├───────────╯ got array, want string
					`,
						filepath.Join(dir, "ConfigRecipeNotString", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigRecipeEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:2:11

						  1 │ manala:
						▶ 2 │   recipe: ""
						    ├───────────╯ minLength: got 0, want 1
					`,
						filepath.Join(dir, "ConfigRecipeEmpty", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigRecipeTooLong",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:2:11

						  1 │ manala:
						▶ 2 │   recipe: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
						    ├───────────╯ maxLength: got 445, want 100
					`,
						filepath.Join(dir, "ConfigRecipeTooLong", "project", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Repository
		{
			test: "ConfigRepositoryNotString",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:15

						  1 │ manala:
						  2 │   recipe: recipe
						▶ 3 │   repository: []
						    ├───────────────╯ got array, want string
					`,
						filepath.Join(dir, "ConfigRepositoryNotString", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigRepositoryEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:15

						  1 │ manala:
						  2 │   recipe: recipe
						▶ 3 │   repository: ""
						    ├───────────────╯ minLength: got 0, want 1
					`,
						filepath.Join(dir, "ConfigRepositoryEmpty", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigRepositoryTooLong",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:15

						  1 │ manala:
						  2 │   recipe: recipe
						▶ 3 │   repository: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
						    ├───────────────╯ maxLength: got 445, want 256
					`,
						filepath.Join(dir, "ConfigRepositoryTooLong", "project", ".manala.yaml"),
					)),
				),
			},
		},
		// Vars
		{
			test: "VarsInvalid",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest vars",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:5:6

						  2 │   recipe: recipe
						  3 │   repository: testdata/LoaderSuite/TestHandleErrors/VarsInvalid/repository
						  4 │
						▶ 5 │ foo: bar
						    ├──────╯ got string, want integer
					`,
						filepath.Join(dir, "VarsInvalid", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "VarsAdditionalProperty",
			expected: serrortest.Expectation{
				Msg: "invalid project manifest vars",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:5:1

						  2 │   recipe: recipe
						  3 │   repository: testdata/LoaderSuite/TestHandleErrors/VarsAdditionalProperty/repository
						  4 │
						▶ 5 │ bar: bar
						    ├─╯ additional property 'bar' not allowed
					`,
						filepath.Join(dir, "VarsAdditionalProperty", "project", ".manala.yaml"),
					)),
				),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			project, err := s.handle(filepath.Join(dir, test.test, "project"))

			s.Nil(project)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}

func (s *LoaderSuite) handle(dir string) (app.Project, error) {
	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
	))
	recipeLoader := recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(
		recipeManifest.NewLoaderHandler(log.Discard),
	))

	chainMock := &project.LoaderHandlerChainMock{}

	handler := manifest.NewLoaderHandler(log.Discard, repositoryLoader, recipeLoader)

	return handler.Handle(&project.LoaderQuery{Dir: dir}, chainMock)
}
