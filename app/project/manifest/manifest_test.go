package manifest_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test() {
	m := manifest.New()

	s.Empty(m.Recipe())
	s.Empty(m.Repository())
	s.Equal(map[string]any{}, m.Vars())
}

func (s *Suite) TestUnmarshalYAMLErrors() {
	tests := []struct {
		test     string
		expected errors.Assertion
	}{
		{
			test: "Empty",
			expected: &parsing.Assertion{
				Err: &serrors.Assertion{
					Message: "empty yaml content",
				},
			},
		},
		{
			test: "Invalid",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "'@' is a reserved character",
				},
			},
		},
		{
			test: "IrregularType",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 6,
				Err: &serrors.Assertion{
					Message: "irregular type",
				},
			},
		},
		{
			test: "IrregularMapKey",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 2,
				Err: &serrors.Assertion{
					Message: "irregular map key",
				},
			},
		},
		{
			test: "NotMap",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "yaml document must be a map",
				},
			},
		},
		// Config
		{
			test: "ConfigAbsent",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "missing manala property",
						Arguments: []any{
							"property", "manala",
							"line", 1,
							"column", 4,
						},
						Details: `
							>  1 | foo: bar
							          ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigNotMap",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "manala field must be a map",
						Arguments: []any{
							"expected", "object",
							"actual", "string",
							"path", "manala",
							"line", 1,
							"column", 9,
						},
						Details: `
							>  1 | manala: foo
							               ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigEmpty",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "missing manala recipe property",
						Arguments: []any{
							"path", "manala",
							"property", "recipe",
							"line", 1,
							"column", 9,
						},
						Details: `
							>  1 | manala: {}
							               ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigAdditionalProperties",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "manala field don't support additional properties",
						Arguments: []any{
							"path", "manala.foo",
							"line", 3,
							"column", 8,
						},
						Details: `
							   1 | manala:
							   2 |   recipe: recipe
							>  3 |   foo: bar
							              ^
						`,
					},
				},
			},
		},
		// Config - Recipe
		{
			test: "ConfigRecipeAbsent",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "missing manala recipe property",
						Arguments: []any{
							"path", "manala",
							"property", "recipe",
							"line", 2,
							"column", 13,
						},
						Details: `
							   1 | manala:
							>  2 |   repository: repository
							                   ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigRecipeNotString",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "manala recipe field must be a string",
						Arguments: []any{
							"expected", "string",
							"actual", "array",
							"path", "manala.recipe",
							"line", 2,
							"column", 11,
						},
						Details: `
							   1 | manala:
							>  2 |   recipe: []
							                 ^
							   3 |   repository: repository
						`,
					},
				},
			},
		},
		{
			test: "ConfigRecipeEmpty",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "empty manala recipe field",
						Arguments: []any{
							"minimum", 1,
							"path", "manala.recipe",
							"line", 2,
							"column", 11,
						},
						Details: `
							   1 | manala:
							>  2 |   recipe: ""
							                 ^
							   3 |   repository: repository
						`,
					},
				},
			},
		},
		{
			test: "ConfigRecipeTooLong",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "too long manala recipe field",
						Arguments: []any{
							"maximum", 100,
							"path", "manala.recipe",
							"line", 2,
							"column", 11,
						},
						Details: `
							   1 | manala:
							>  2 |   recipe: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							                 ^
							   3 |   repository: repository
						`,
					},
				},
			},
		},
		// Config - Repository
		{
			test: "ConfigRepositoryNotString",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "manala repository field must be a string",
						Arguments: []any{
							"expected", "string",
							"actual", "array",
							"path", "manala.repository",
							"line", 3,
							"column", 15,
						},
						Details: `
							   1 | manala:
							   2 |   recipe: recipe
							>  3 |   repository: []
							                     ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigRepositoryEmpty",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "empty manala repository field",
						Arguments: []any{
							"minimum", 1,
							"path", "manala.repository",
							"line", 3,
							"column", 15,
						},
						Details: `
							   1 | manala:
							   2 |   recipe: recipe
							>  3 |   repository: ""
							                     ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigRepositoryTooLong",
			expected: &serrors.Assertion{
				Message: "invalid project manifest",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "too long manala repository field",
						Arguments: []any{
							"maximum", 256,
							"path", "manala.repository",
							"line", 3,
							"column", 15,
						},
						Details: `
							   1 | manala:
							   2 |   recipe: recipe
							>  3 |   repository: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							                     ^
						`,
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			m := manifest.New()

			dir := filepath.FromSlash("testdata/Suite/TestUnmarshalYAMLErrors")

			reader, _ := os.Open(filepath.Join(dir, test.test+".yaml"))
			content, _ := io.ReadAll(reader)

			err := m.UnmarshalYAML(content)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *Suite) TestUnmarshalYAML() {
	tests := []struct {
		test               string
		expectedRecipe     string
		expectedRepository string
		expectedVars       map[string]any
	}{
		{
			test:               "All",
			expectedRecipe:     "recipe",
			expectedRepository: "repository",
			expectedVars: map[string]any{
				"foo": "bar",
			},
		},
		{
			test:               "ConfigRepositoryAbsent",
			expectedRecipe:     "recipe",
			expectedRepository: "",
			expectedVars: map[string]any{
				"foo": "bar",
			},
		},
		{
			test:               "VarsAbsent",
			expectedRecipe:     "recipe",
			expectedRepository: "repository",
			expectedVars:       map[string]any{},
		},
		{
			test:               "VarsKeys",
			expectedRecipe:     "recipe",
			expectedRepository: "repository",
			expectedVars: map[string]any{
				"underscore_key": "ok",
				"hyphen-key":     "ok",
				"dot.key":        "ok",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			m := manifest.New()

			dir := filepath.FromSlash("testdata/Suite/TestUnmarshalYAML")

			reader, _ := os.Open(filepath.Join(dir, test.test+".yaml"))
			content, _ := io.ReadAll(reader)

			err := m.UnmarshalYAML(content)

			s.Require().NoError(err)
			s.Equal(test.expectedRecipe, m.Recipe())
			s.Equal(test.expectedRepository, m.Repository())
			s.Equal(test.expectedVars, m.Vars())
		})
	}
}
