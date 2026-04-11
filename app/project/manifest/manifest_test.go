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
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "empty yaml content",
				},
			},
		},
		{
			test: "Invalid",
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
			expected: &parsing.ErrorAssertion{
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
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "missing manala property",
				},
			},
		},
		{
			test: "ConfigNotMap",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "string was used where mapping is expected",
				},
			},
		},
		{
			test: "ConfigEmpty",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "missing manala recipe property",
				},
			},
		},
		{
			test: "ConfigAdditionalProperties",
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 3,
				Err: &serrors.Assertion{
					Message: "unknown field \"foo\"",
				},
			},
		},
		// Config - Recipe
		{
			test: "ConfigRecipeAbsent",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 7,
				Err: &serrors.Assertion{
					Message: "missing manala recipe property",
				},
			},
		},
		{
			test: "ConfigRecipeNotString",
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigRecipeEmpty",
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "missing manala recipe property",
				},
			},
		},
		{
			test: "ConfigRecipeTooLong",
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "too long manala recipe field (max=100)",
				},
			},
		},
		// Config - Repository
		{
			test: "ConfigRepositoryNotString",
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 15,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigRepositoryTooLong",
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
