package manifest_test

import (
	"testing"

	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/testing/heredoc"
	"github.com/manala/manala/internal/validation"
	yamlparser "github.com/manala/manala/internal/yaml/parser"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct{ suite.Suite }

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) TestUnmarshalRequired() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		recipe: recipe
	`)))
	s.Require().NoError(err)

	c := &manifest.Config{}
	err = c.UnmarshalYAML(node)
	s.Require().NoError(err)

	s.Equal("recipe", c.Recipe)
	s.Empty(c.Repository)
}

func (s *ConfigSuite) TestUnmarshal() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		recipe: recipe
		repository: repository
	`)))
	s.Require().NoError(err)

	c := &manifest.Config{}
	err = c.UnmarshalYAML(node)
	s.Require().NoError(err)

	s.Equal("recipe", c.Recipe)
	s.Equal("repository", c.Repository)
}

func (s *ConfigSuite) TestUnmarshalErrors() {
	tests := []struct {
		test     string
		content  string
		expected expectation.ErrorExpectation
	}{
		{
			test: "Empty",
			content: heredoc.Doc(`
				{}
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 1},
					Err:      expectation.ErrorMessage("missing property 'recipe'"),
				},
			),
		},
		{
			test: "AdditionalProperties",
			content: heredoc.Doc(`
				recipe: recipe
				foo: bar
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 7},
					Err:      expectation.ErrorMessage("additional properties 'foo' not allowed"),
				},
			),
		},
		// Recipe
		{
			test: "RecipeAbsent",
			content: heredoc.Doc(`
				repository: repository
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("missing property 'recipe'"),
				},
			),
		},
		{
			test: "RecipeNotString",
			content: heredoc.Doc(`
				recipe: []
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 9},
					Err:      expectation.ErrorMessage("got array, want string"),
				},
			),
		},
		{
			test: "RecipeEmpty",
			content: heredoc.Doc(`
				recipe: ""
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 9},
					Err:      expectation.ErrorMessage("minLength: got 0, want 1"),
				},
			),
		},
		{
			test: "RecipeTooLong",
			content: heredoc.Doc(`
				recipe: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 9},
					Err:      expectation.ErrorMessage("maxLength: got 445, want 100"),
				},
			),
		},
		// Repository
		{
			test: "RepositoryNotString",
			content: heredoc.Doc(`
				recipe: recipe
				repository: []
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{2, 13},
					Err:      expectation.ErrorMessage("got array, want string"),
				},
			),
		},
		{
			test: "RepositoryTooLong",
			content: heredoc.Doc(`
				recipe: recipe
				repository: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{2, 13},
					Err:      expectation.ErrorMessage("maxLength: got 445, want 256"),
				},
			),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := yamlparser.Parse([]byte(test.content))
			s.Require().NoError(err)

			c := &manifest.Config{}
			err = c.UnmarshalYAML(node)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
