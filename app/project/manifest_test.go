package project

import (
	"github.com/stretchr/testify/suite"
	"manala/internal/serrors"
	"os"
	"path/filepath"
	"testing"
)

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) Test() {
	manifest := NewManifest()

	s.Equal("", manifest.Recipe())
	s.Equal("", manifest.Repository())
	s.Equal(map[string]any{}, manifest.Vars())
}

func (s *ManifestSuite) TestReadFromErrors() {
	tests := []struct {
		test     string
		expected *serrors.Assert
	}{
		{
			test: "Empty",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "empty yaml file",
					},
				},
			},
		},
		{
			test: "Invalid",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "unexpected mapping key",
						Arguments: []any{
							"line", 1,
							"column", 1,
						},
						Details: `
							>  1 | ::
							       ^
						`,
					},
				},
			},
		},
		{
			test: "IrregularType",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "irregular type",
						Arguments: []any{
							"line", 1,
							"column", 6,
						},
						Details: `
							>  1 | foo: .inf
							            ^
						`,
					},
				},
			},
		},
		{
			test: "IrregularMapKey",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "irregular map key",
						Arguments: []any{
							"line", 1,
							"column", 2,
						},
						Details: `
							>  1 | 0: bar
							        ^
						`,
					},
				},
			},
		},
		{
			test: "NotMap",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "yaml document must be a map",
						Arguments: []any{
							"expected", "object",
							"actual", "string",
						},
					},
				},
			},
		},
		// Config
		{
			test: "ConfigAbsent",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "missing manala property",
						Arguments: []any{
							"property", "manala",
						},
					},
				},
			},
		},
		{
			test: "ConfigNotMap",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			manifest := NewManifest()

			dir := filepath.FromSlash("testdata/ManifestSuite/TestReadFromErrors/" + test.test)

			manifestFile, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
			_, err := manifest.ReadFrom(manifestFile)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *ManifestSuite) TestReadFrom() {
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
			manifest := NewManifest()

			dir := filepath.FromSlash("testdata/ManifestSuite/TestReadFrom/" + test.test)

			manifestFile, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
			_, err := manifest.ReadFrom(manifestFile)

			s.NoError(err)

			s.Equal(test.expectedRecipe, manifest.Recipe())
			s.Equal(test.expectedRepository, manifest.Repository())
			s.Equal(test.expectedVars, manifest.Vars())
		})
	}
}
