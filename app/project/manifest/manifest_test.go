package manifest_test

import (
	"testing"

	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) TestNew() {
	m := manifest.New()

	s.Empty(m.Recipe)
	s.Empty(m.Repository)
	s.Equal(map[string]any{}, m.Vars())
}

func (s *ManifestSuite) TestUnmarshalRequired() {
	m := manifest.New()

	err := m.Unmarshal([]byte(heredoc.Doc(`
		manala:
		  recipe: recipe
	`)))
	s.Require().NoError(err)

	s.Equal("recipe", m.Recipe)
	s.Empty(m.Repository)
	s.Equal(map[string]any{}, m.Vars())
}

func (s *ManifestSuite) TestUnmarshal() {
	m := manifest.New()

	err := m.Unmarshal([]byte(heredoc.Doc(`
		manala:
		  recipe: recipe
		  repository: repository
		foo: bar
		underscore_key: ok
		hyphen-key: ok
		dot.key: ok
	`)))
	s.Require().NoError(err)

	s.Equal("recipe", m.Recipe)
	s.Equal("repository", m.Repository)
	s.Equal(map[string]any{
		"foo":            "bar",
		"underscore_key": "ok",
		"hyphen-key":     "ok",
		"dot.key":        "ok",
	}, m.Vars())
}

func (s *ManifestSuite) TestUnmarshalErrors() {
	tests := []struct {
		test     string
		content  string
		expected errors.Assertion
	}{
		{
			test:    "Empty",
			content: "",
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "empty yaml content",
				},
			},
		},
		{
			test: "Invalid",
			content: heredoc.Doc(`
				@
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "'@' is a reserved character",
				},
			},
		},
		{
			test: "IrregularType",
			content: heredoc.Doc(`
				foo: .inf
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 6,
				Err: &serrors.Assertion{
					Message: "irregular type",
				},
			},
		},
		{
			test: "IrregularMapKey",
			content: heredoc.Doc(`
				0: foo
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 2,
				Err: &serrors.Assertion{
					Message: "irregular map key",
				},
			},
		},
		{
			test: "NotMap",
			content: heredoc.Doc(`
				foo
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "yaml document must be a map",
				},
			},
		},
		// Manala
		{
			test: "ManalaAbsent",
			content: heredoc.Doc(`
				foo: bar
			`),
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "missing manala property",
				},
			},
		},
		{
			test: "ManalaNotMap",
			content: heredoc.Doc(`
				manala: foo
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "string was used where mapping is expected",
				},
			},
		},
		{
			test: "ManalaEmpty",
			content: heredoc.Doc(`
				manala: {}
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "missing manala recipe property",
				},
			},
		},
		{
			test: "ManalaAdditionalProperties",
			content: heredoc.Doc(`
				manala:
				  recipe: recipe
				  foo: bar
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 3,
				Err: &serrors.Assertion{
					Message: "unknown field \"foo\"",
				},
			},
		},
		// Recipe
		{
			test: "RecipeAbsent",
			content: heredoc.Doc(`
				manala:
				  repository: repository
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 7,
				Err: &serrors.Assertion{
					Message: "missing manala recipe property",
				},
			},
		},
		{
			test: "RecipeNotString",
			content: heredoc.Doc(`
				manala:
				  recipe: []
				  repository: repository
			`),
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "RecipeEmpty",
			content: heredoc.Doc(`
				manala:
				  recipe: ""
				  repository: repository
			`),
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "missing manala recipe property",
				},
			},
		},
		{
			test: "RecipeTooLong",
			content: heredoc.Doc(`
				manala:
				  recipe: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
				  repository: repository
			`),
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "too long manala recipe field (max=100)",
				},
			},
		},
		// Repository
		{
			test: "RepositoryNotString",
			content: heredoc.Doc(`
				manala:
				  recipe: recipe
				  repository: []
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 15,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "RepositoryTooLong",
			content: heredoc.Doc(`
				manala:
				  recipe: recipe
				  repository: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 15,
				Err: &serrors.Assertion{
					Message: "too long manala repository field (max=256)",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			m := manifest.New()

			err := m.Unmarshal([]byte(test.content))

			errors.Equal(s.T(), test.expected, err)
		})
	}
}
