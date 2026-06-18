package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/app/sync"
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
	repositoryURL := filepath.FromSlash("testdata/LoaderSuite/TestHandle/repository")

	recipe, err := s.handle(repositoryURL)

	s.Require().NoError(err)

	s.Equal(filepath.Join(repositoryURL, "recipe"), recipe.Dir())
	s.Equal("recipe", recipe.Name())
	s.Equal("description", recipe.Description())
	s.Equal("icon", recipe.Icon())
	s.Equal(filepath.Join(repositoryURL, "recipe", "template"), recipe.Template())
	s.Equal([]string{
		filepath.Join(repositoryURL, "recipe", "partial.tmpl"),
		filepath.Join(repositoryURL, "recipe", "dir", "partial.tmpl"),
	}, recipe.Partials())
	sync.ExpectUnits(s.T(), sync.UnitExpectations{
		{Source: "file", Destination: "file"},
		{Source: "dir/file", Destination: "dir/file"},
		{Source: "file", Destination: "dir/file"},
		{Source: "dir/file", Destination: "file"},
		{Source: "src_file", Destination: "dst_file"},
		{Source: "src_dir/file", Destination: "dst_dir/file"},
	}, recipe.Sync())
	s.Equal(repositoryURL, recipe.Repository().URL())
	s.Equal(map[string]any{"foo": nil, "bar": "baz"}, recipe.Vars())
	s.Equal(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"foo": map[string]any{"type": "int"},
			"bar": map[string]any{"type": "string"},
		},
		"additionalProperties": false,
	}, recipe.Schema())
	option.ExpectOptions(s.T(), option.Expectations{{
		Type:      &option.String{},
		Label:     "label",
		Name:      "name",
		MaxLength: 0,
		Values:    []any{},
	}}, recipe.Options())
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
				Msg: "recipe manifest is a directory",
				Attrs: [][2]any{
					{"dir", filepath.Join(dir, "Directory", "repository", "recipe", ".manala.yaml")},
				},
			},
		},
		{
			test: "Unparsable",
			expected: serrortest.Expectation{
				Msg: "unable to parse recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

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
			test: "Empty",
			expected: serrortest.Expectation{
				Msg: "unable to parse recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s

						  1 │

						empty yaml content
					`,
						filepath.Join(dir, "Empty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "MultipleDocuments",
			expected: serrortest.Expectation{
				Msg: "unable to parse recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:2:1

						  1 │ foo: bar
						▶ 2 │ ---
						    ├─╯ multiple documents yaml content
						  3 │ foo: bar
					`,
						filepath.Join(dir, "MultipleDocuments", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "NotMap",
			expected: serrortest.Expectation{
				Msg: "unable to parse recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:1:1

						▶ 1 │ - foo
						    ├─╯ yaml document must be a map
						  2 │ - bar
					`,
						filepath.Join(dir, "NotMap", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "MapEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s

						  1 │ {}

						missing property 'manala'
					`,
						filepath.Join(dir, "MapEmpty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		// Config
		{
			test: "ConfigMissing",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s

						  1 │ foo: bar

						missing property 'manala'
					`,
						filepath.Join(dir, "ConfigMissing", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigNotMap",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:1:9

						▶ 1 │ manala: foo
						    ├─────────╯ got string, want object
					`,
						filepath.Join(dir, "ConfigNotMap", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigMapEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:1:1

						▶ 1 │ manala: {}
						    ├─╯ missing property 'description'
					`,
						filepath.Join(dir, "ConfigMapEmpty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigAdditionalProperty",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:3

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   foo: bar
						    ├───╯ additional property 'foo' not allowed
					`,
						filepath.Join(dir, "ConfigAdditionalProperty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Description
		{
			test: "ConfigDescriptionAbsent",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:1:1

						▶ 1 │ manala:
						    ├─╯ missing property 'description'
						  2 │   icon: icon
					`,
						filepath.Join(dir, "ConfigDescriptionAbsent", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigDescriptionNotString",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:2:16

						  1 │ manala:
						▶ 2 │   description: []
						    ├────────────────╯ got array, want string
					`,
						filepath.Join(dir, "ConfigDescriptionNotString", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigDescriptionEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:2:16

						  1 │ manala:
						▶ 2 │   description: ""
						    ├────────────────╯ minLength: got 0, want 1
					`,
						filepath.Join(dir, "ConfigDescriptionEmpty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigDescriptionTooLong",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:2:16

						  1 │ manala:
						▶ 2 │   description: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
						    ├────────────────╯ maxLength: got 445, want 256
					`,
						filepath.Join(dir, "ConfigDescriptionTooLong", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Icon
		{
			test: "ConfigIconNotString",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:9

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   icon: []
						    ├─────────╯ got array, want string
					`,
						filepath.Join(dir, "ConfigIconNotString", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigIconEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:9

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   icon: ""
						    ├─────────╯ minLength: got 0, want 1
					`,
						filepath.Join(dir, "ConfigIconEmpty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigIconTooLong",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:9

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   icon: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
						    ├─────────╯ maxLength: got 445, want 100
					`,
						filepath.Join(dir, "ConfigIconTooLong", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Template
		{
			test: "ConfigTemplateNotString",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:13

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   template: []
						    ├─────────────╯ got array, want string
					`,
						filepath.Join(dir, "ConfigTemplateNotString", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigTemplateEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:13

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   template: ""
						    ├─────────────╯ minLength: got 0, want 1
					`,
						filepath.Join(dir, "ConfigTemplateEmpty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigTemplateTooLong",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:13

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   template: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
						    ├─────────────╯ maxLength: got 445, want 100
					`,
						filepath.Join(dir, "ConfigTemplateTooLong", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Partials
		{
			test: "ConfigPartialsNotArray",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:13

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   partials: foo
						    ├─────────────╯ got string, want array
					`,
						filepath.Join(dir, "ConfigPartialsNotArray", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Partials Item
		{
			test: "ConfigPartialsItemNotString",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:4:7

						  1 │ manala:
						  2 │   description: description
						  3 │   partials:
						▶ 4 │     - []
						    ├───────╯ got array, want string
					`,
						filepath.Join(dir, "ConfigPartialsItemNotString", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigPartialsItemEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:4:7

						  1 │ manala:
						  2 │   description: description
						  3 │   partials:
						▶ 4 │     - ""
						    ├───────╯ minLength: got 0, want 1
					`,
						filepath.Join(dir, "ConfigPartialsItemEmpty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigPartialsItemTooLong",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:4:7

						  1 │ manala:
						  2 │   description: description
						  3 │   partials:
						▶ 4 │     - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
						    ├───────╯ maxLength: got 445, want 100
					`,
						filepath.Join(dir, "ConfigPartialsItemTooLong", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Sync
		{
			test: "ConfigSyncNotArray",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:3:9

						  1 │ manala:
						  2 │   description: description
						▶ 3 │   sync: foo
						    ├─────────╯ got string, want array
					`,
						filepath.Join(dir, "ConfigSyncNotArray", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		// Config - Sync Item
		{
			test: "ConfigSyncItemNotString",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:4:7

						  1 │ manala:
						  2 │   description: description
						  3 │   sync:
						▶ 4 │     - []
						    ├───────╯ got array, want string
					`,
						filepath.Join(dir, "ConfigSyncItemNotString", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigSyncItemEmpty",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:4:7

						  1 │ manala:
						  2 │   description: description
						  3 │   sync:
						▶ 4 │     - ""
						    ├───────╯ minLength: got 0, want 1
					`,
						filepath.Join(dir, "ConfigSyncItemEmpty", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "ConfigSyncItemTooLong",
			expected: serrortest.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:4:7

						  1 │ manala:
						  2 │   description: description
						  3 │   sync:
						▶ 4 │     - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
						    ├───────╯ maxLength: got 445, want 256
					`,
						filepath.Join(dir, "ConfigSyncItemTooLong", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "AnnotationUnparsableSingleLine",
			expected: serrortest.Expectation{
				Msg: "unable to infer recipe manifest vars",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:5:14

						  2 │   description: description
						  3 │
						  4 │ foo:
						▶ 5 │   # @schema foo
						    ├──────────────╯ invalid character 'o' in literal false (expecting 'a')
						  6 │   bar: ~
					`,
						filepath.Join(dir, "AnnotationUnparsableSingleLine", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "AnnotationUnparsableMultiLine",
			expected: serrortest.Expectation{
				Msg: "unable to infer recipe manifest vars",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:6:8

						  4 │ foo:
						  5 │   # @schema
						▶ 6 │   #   foo
						    ├────────╯ invalid character 'o' in literal false (expecting 'a')
						  7 │   bar: ~
					`,
						filepath.Join(dir, "AnnotationUnparsableMultiLine", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "AnnotationInvalidSingleLine",
			expected: serrortest.Expectation{
				Msg: "unable to infer recipe manifest vars",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:4:30

						  1 │ manala:
						  2 │   description: description
						  3 │
						▶ 4 │ # @option {"label": "Label", "foo": "bar"}
						    ├──────────────────────────────╯ additional property 'foo' not allowed
						  5 │ node: foo
					`,
						filepath.Join(dir, "AnnotationInvalidSingleLine", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "AnnotationInvalidMultiLine",
			expected: serrortest.Expectation{
				Msg: "unable to infer recipe manifest vars",
				Err: expectation.Errors(
					sourcetest.Expectation(heredoc.Doc(`

						at %[1]s:5:5

						  2 │   description: description
						  3 │
						  4 │ # @option {"label": "Label",
						▶ 5 │ #   "foo": "bar"
						    ├─────╯ additional property 'foo' not allowed
						  6 │ # }
						  7 │ node: foo
					`,
						filepath.Join(dir, "AnnotationInvalidMultiLine", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			recipe, err := s.handle(filepath.Join(dir, test.test, "repository"))

			s.Nil(recipe)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}

func (s *LoaderSuite) handle(repositoryURL string) (app.Recipe, error) {
	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(log.Discard),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	chainMock := &recipe.LoaderHandlerChainMock{}

	handler := manifest.NewLoaderHandler(log.Discard)
	return handler.Handle(&recipe.LoaderQuery{Repository: repository, Name: "recipe"}, chainMock)
}
