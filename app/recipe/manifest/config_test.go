package manifest_test

import (
	"testing"

	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/sync"
	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/testing/heredoc"
	"github.com/manala/manala/internal/validation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlparser "github.com/manala/manala/internal/yaml/parser"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct{ suite.Suite }

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) TestUnmarshalRequired() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		description: description
	`)))
	s.Require().NoError(err)

	c := &manifest.Config{}
	err = c.UnmarshalYAML(node)
	s.Require().NoError(err)

	s.Equal("description", c.Description)
	s.Empty(c.Icon)
	s.Empty(c.Template)
	s.Empty(c.Partials)
	sync.ExpectUnits(s.T(), sync.UnitsExpectation{}, c.Sync)
}

func (s *ConfigSuite) TestUnmarshal() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		description: description
		icon: icon
		template: template
		partials:
		  - partial.tmpl
		  - dir/partial.tmpl
		sync:
		  - file
		  - dir/file
		  - file dir/file
		  - dir/file file
		  - src_file dst_file
		  - src_dir/file dst_dir/file
	`)))
	s.Require().NoError(err)

	c := &manifest.Config{}
	err = c.UnmarshalYAML(node)
	s.Require().NoError(err)

	s.Equal("description", c.Description)
	s.Equal("icon", c.Icon)
	s.Equal("template", c.Template)
	s.Equal([]string{
		"partial.tmpl",
		"dir/partial.tmpl",
	}, c.Partials)
	sync.ExpectUnits(s.T(), sync.UnitsExpectation{
		{Source: "file", Destination: "file"},
		{Source: "dir/file", Destination: "dir/file"},
		{Source: "file", Destination: "dir/file"},
		{Source: "dir/file", Destination: "file"},
		{Source: "src_file", Destination: "dst_file"},
		{Source: "src_dir/file", Destination: "dst_dir/file"},
	}, c.Sync)
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
					Err:      expectation.ErrorMessage("missing property 'description'"),
				},
			),
		},
		{
			test: "AdditionalProperties",
			content: heredoc.Doc(`
				description: description
				foo: bar
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 12},
					Err:      expectation.ErrorMessage("additional properties 'foo' not allowed"),
				},
			),
		},
		// Description
		{
			test: "DescriptionAbsent",
			content: heredoc.Doc(`
				template: template
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 9},
					Err:      expectation.ErrorMessage("missing property 'description'"),
				},
			),
		},
		{
			test: "DescriptionNotString",
			content: heredoc.Doc(`
				description: []
				template: template
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 14},
					Err:      expectation.ErrorMessage("got array, want string"),
				},
			),
		},
		{
			test: "DescriptionEmpty",
			content: heredoc.Doc(`
				description: ""
				template: template
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 14},
					Err:      expectation.ErrorMessage("minLength: got 0, want 1"),
				},
			),
		},
		{
			test: "DescriptionTooLong",
			content: heredoc.Doc(`
				description: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
				template: template
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{1, 14},
					Err:      expectation.ErrorMessage("maxLength: got 445, want 256"),
				},
			),
		},
		// Icon
		{
			test: "IconNotString",
			content: heredoc.Doc(`
				description: description
				icon: []
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{2, 7},
					Err:      expectation.ErrorMessage("got array, want string"),
				},
			),
		},
		{
			test: "IconTooLong",
			content: heredoc.Doc(`
				description: description
				icon: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{2, 7},
					Err:      expectation.ErrorMessage("maxLength: got 445, want 100"),
				},
			),
		},
		// Template
		{
			test: "TemplateNotString",
			content: heredoc.Doc(`
				description: description
				template: []
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{2, 11},
					Err:      expectation.ErrorMessage("got array, want string"),
				},
			),
		},
		{
			test: "TemplateTooLong",
			content: heredoc.Doc(`
				description: description
				template: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{2, 11},
					Err:      expectation.ErrorMessage("maxLength: got 445, want 100"),
				},
			),
		},
		// Partials
		{
			test: "PartialsNotArray",
			content: heredoc.Doc(`
				description: description
				partials: foo
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{2, 11},
					Err:      expectation.ErrorMessage("got string, want array"),
				},
			),
		},
		// Partials Item
		{
			test: "PartialsItemNotString",
			content: heredoc.Doc(`
				description: description
				partials:
				  - []
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{3, 5},
					Err:      expectation.ErrorMessage("got array, want string"),
				},
			),
		},
		{
			test: "PartialsItemEmpty",
			content: heredoc.Doc(`
				description: description
				partials:
				  - ""
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{3, 5},
					Err:      expectation.ErrorMessage("minLength: got 0, want 1"),
				},
			),
		},
		{
			test: "PartialsItemTooLong",
			content: heredoc.Doc(`
				description: description
				partials:
				  - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{3, 5},
					Err:      expectation.ErrorMessage("maxLength: got 445, want 100"),
				},
			),
		},
		// Sync
		{
			test: "SyncNotArray",
			content: heredoc.Doc(`
				description: description
				sync: foo
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{2, 7},
				Err:      expectation.ErrorMessage("string was used where sequence is expected"),
			},
		},
		// Sync Item
		{
			test: "SyncItemNotString",
			content: heredoc.Doc(`
				description: description
				sync:
				  - []
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{3, 5},
				Err:      expectation.ErrorMessage("field must be a string"),
			},
		},
		{
			test: "SyncItemEmpty",
			content: heredoc.Doc(`
				description: description
				sync:
				  - ""
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{3, 5},
					Err:      expectation.ErrorMessage("minLength: got 0, want 1"),
				},
			),
		},
		{
			test: "SyncItemTooLong",
			content: heredoc.Doc(`
				description: description
				sync:
				  - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: expectation.Errors(
				validation.ViolationExpectation{
					Position: [2]int{3, 5},
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
