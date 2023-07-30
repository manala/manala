package project

import (
	"github.com/stretchr/testify/suite"
	"manala/internal/errors/serrors"
	"manala/internal/testing/heredoc"
	"manala/internal/validation"
	"manala/internal/yaml"
	"os"
	"path/filepath"
	"testing"
)

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) Test() {
	recMan := NewManifest()

	s.Equal("", recMan.Recipe())
	s.Equal("", recMan.Repository())
	s.Equal(map[string]interface{}{}, recMan.Vars())
}

func (s *ManifestSuite) TestReadFromErrors() {
	tests := []struct {
		test     string
		expected *serrors.Assert
	}{
		{
			test: "Empty",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular project manifest",
				Error: &serrors.Assert{
					Type:    &serrors.Error{},
					Message: "empty yaml file",
				},
			},
		},
		{
			test: "Invalid",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular project manifest",
				Error: &serrors.Assert{
					Type:    &yaml.Error{},
					Message: "unexpected mapping key",
					Arguments: []any{
						"line", 1,
						"column", 1,
					},
					Details: heredoc.Doc(`
						>  1 | ::
						       ^
					`),
				},
			},
		},
		{
			test: "IrregularType",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular project manifest",
				Error: &serrors.Assert{
					Type:    &yaml.NodeError{},
					Message: "irregular type",
					Arguments: []any{
						"line", 1,
						"column", 6,
					},
					Details: heredoc.Doc(`
						>  1 | foo: .inf
						            ^
					`),
				},
			},
		},
		{
			test: "IrregularMapKey",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular project manifest",
				Error: &serrors.Assert{
					Type:    &yaml.NodeError{},
					Message: "irregular map key",
					Arguments: []any{
						"line", 1,
						"column", 2,
					},
					Details: heredoc.Doc(`
						>  1 | 0: bar
						        ^
					`),
				},
			},
		},
		{
			test: "NotMap",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "yaml document must be a map",
						Arguments: []any{
							"expected", "object",
							"given", "string",
						},
					},
				},
			},
		},
		// Config
		{
			test: "ConfigAbsent",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala field",
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
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala field must be a map",
						Arguments: []any{
							"expected", "object",
							"given", "string",
							"line", 1,
							"column", 9,
						},
						Details: heredoc.Doc(`
							>  1 | manala: foo
							               ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigEmpty",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala recipe field",
						Arguments: []any{
							"property", "recipe",
							"line", 1,
							"column", 9,
						},
						Details: heredoc.Doc(`
							>  1 | manala: {}
							               ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigAdditionalProperties",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala field don't support additional properties",
						Arguments: []any{
							"property", "foo",
							"line", 2,
							"column", 9,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   recipe: recipe
							               ^
							   3 |   foo: bar
						`),
					},
				},
			},
		},
		// Config - Recipe
		{
			test: "ConfigRecipeAbsent",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala recipe field",
						Arguments: []any{
							"property", "recipe",
							"line", 2,
							"column", 13,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   repository: repository
							                   ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigRecipeNotString",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala recipe field must be a string",
						Arguments: []any{
							"expected", "string",
							"given", "array",
							"line", 2,
							"column", 11,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   recipe: []
							                 ^
							   3 |   repository: repository
						`),
					},
				},
			},
		},
		{
			test: "ConfigRecipeEmpty",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "empty manala recipe field",
						Arguments: []any{
							"line", 2,
							"column", 11,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   recipe: ""
							                 ^
							   3 |   repository: repository
						`),
					},
				},
			},
		},
		{
			test: "ConfigRecipeTooLong",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "too long manala recipe field",
						Arguments: []any{
							"line", 2,
							"column", 11,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   recipe: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							                 ^
							   3 |   repository: repository
						`),
					},
				},
			},
		},
		// Config - Repository
		{
			test: "ConfigRepositoryNotString",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala repository field must be a string",
						Arguments: []any{
							"expected", "string",
							"given", "array",
							"line", 3,
							"column", 15,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   recipe: recipe
							>  3 |   repository: []
							                     ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigRepositoryEmpty",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "empty manala repository field",
						Arguments: []any{
							"line", 3,
							"column", 15,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   recipe: recipe
							>  3 |   repository: ""
							                     ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigRepositoryTooLong",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "too long manala repository field",
						Arguments: []any{
							"line", 3,
							"column", 15,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   recipe: recipe
							>  3 |   repository: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							                     ^
						`),
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			projMan := NewManifest()

			projDir := filepath.FromSlash("testdata/ManifestSuite/TestReadFromErrors/" + test.test)

			projManFile, _ := os.Open(filepath.Join(projDir, "manifest.yaml"))
			err := projMan.ReadFrom(projManFile)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *ManifestSuite) TestReadFrom() {
	tests := []struct {
		test               string
		expectedRecipe     string
		expectedRepository string
		expectedVars       map[string]interface{}
	}{
		{
			test:               "All",
			expectedRecipe:     "recipe",
			expectedRepository: "repository",
			expectedVars: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			test:               "ConfigRepositoryAbsent",
			expectedRecipe:     "recipe",
			expectedRepository: "",
			expectedVars: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			test:               "VarsAbsent",
			expectedRecipe:     "recipe",
			expectedRepository: "repository",
			expectedVars:       map[string]interface{}{},
		},
		{
			test:               "VarsKeys",
			expectedRecipe:     "recipe",
			expectedRepository: "repository",
			expectedVars: map[string]interface{}{
				"underscore_key": "ok",
				"hyphen-key":     "ok",
				"dot.key":        "ok",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			projMan := NewManifest()

			projDir := filepath.FromSlash("testdata/ManifestSuite/TestReadFrom/" + test.test)

			projManFile, _ := os.Open(filepath.Join(projDir, "manifest.yaml"))
			err := projMan.ReadFrom(projManFile)

			s.NoError(err)

			s.Equal(test.expectedRecipe, projMan.Recipe())
			s.Equal(test.expectedRepository, projMan.Repository())
			s.Equal(test.expectedVars, projMan.Vars())
		})
	}
}
