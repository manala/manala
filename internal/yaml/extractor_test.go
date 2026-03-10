package yaml_test

import (
	"path/filepath"
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/yaml"

	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type ExtractorSuite struct{ suite.Suite }

func TestExtractorSuite(t *testing.T) {
	suite.Run(t, new(ExtractorSuite))
}

func (s *ExtractorSuite) TestRootMapErrors() {
	tests := []struct {
		test     string
		expected *serrors.Assertion
	}{
		{
			test: "Empty",
			expected: &serrors.Assertion{
				Message: "root must be a map",
			},
		},
		{
			test: "NonMap",
			expected: &serrors.Assertion{
				Message: "root must be a map",
				Arguments: []any{
					"line", 1,
					"column", 1,
				},
				Details: `
					>  1 | foo
					       ^
				`,
			},
		},
		{
			test: "SubjectNotFoundSingle",
			expected: &serrors.Assertion{
				Message: "unable to find map",
				Arguments: []any{
					"key", "subject",
				},
			},
		},
		{
			test: "SubjectNotFoundMultiple",
			expected: &serrors.Assertion{
				Message: "unable to find map",
				Arguments: []any{
					"key", "subject",
				},
			},
		},
		{
			test: "SubjectNonMapSingle",
			expected: &serrors.Assertion{
				Message: "key is not a map",
				Arguments: []any{
					"line", 1,
					"column", 10,
					"key", "subject",
				},
				Details: `
					>  1 | subject: 123
					                ^
				`,
			},
		},
		{
			test: "SubjectNonMapMultiple",
			expected: &serrors.Assertion{
				Message: "key is not a map",
				Arguments: []any{
					"line", 1,
					"column", 10,
					"key", "subject",
				},
				Details: `
					>  1 | subject: 123
					                ^
					   2 | foo: foo
				`,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			dir := filepath.FromSlash("testdata/ExtractorSuite/TestRootMapErrors/" + test.test)

			parser := yaml.NewParser()
			node, _ := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			extractor := yaml.NewExtractor(&node)
			subjectNode, err := extractor.ExtractRootMap("subject")

			s.Nil(subjectNode)

			serrors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *ExtractorSuite) TestRootMap() {
	tests := []struct {
		test            string
		expectedSubject any
		expectedNode    any
	}{
		{
			test:            "SingleEmpty",
			expectedSubject: (*goYamlAst.MappingNode)(nil),
			expectedNode:    (*goYamlAst.MappingNode)(nil),
		},
		{
			test:            "SingleSingle",
			expectedSubject: (*goYamlAst.MappingValueNode)(nil),
			expectedNode:    (*goYamlAst.MappingNode)(nil),
		},
		{
			test:            "SingleMultiple",
			expectedSubject: (*goYamlAst.MappingNode)(nil),
			expectedNode:    (*goYamlAst.MappingNode)(nil),
		},
		{
			test:            "CoupleEmpty",
			expectedSubject: (*goYamlAst.MappingNode)(nil),
			expectedNode:    (*goYamlAst.MappingValueNode)(nil),
		},
		{
			test:            "CoupleSingle",
			expectedSubject: (*goYamlAst.MappingValueNode)(nil),
			expectedNode:    (*goYamlAst.MappingValueNode)(nil),
		},
		{
			test:            "CoupleMultiple",
			expectedSubject: (*goYamlAst.MappingNode)(nil),
			expectedNode:    (*goYamlAst.MappingValueNode)(nil),
		},
		{
			test:            "MultipleEmpty",
			expectedSubject: (*goYamlAst.MappingNode)(nil),
			expectedNode:    (*goYamlAst.MappingNode)(nil),
		},
		{
			test:            "MultipleSingle",
			expectedSubject: (*goYamlAst.MappingValueNode)(nil),
			expectedNode:    (*goYamlAst.MappingNode)(nil),
		},
		{
			test:            "MultipleMultiple",
			expectedSubject: (*goYamlAst.MappingNode)(nil),
			expectedNode:    (*goYamlAst.MappingNode)(nil),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			dir := filepath.FromSlash("testdata/ExtractorSuite/TestRootMap/" + test.test)

			parser := yaml.NewParser()
			node, _ := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			extractor := yaml.NewExtractor(&node)
			subjectNode, err := extractor.ExtractRootMap("subject")

			s.Require().NoError(err)

			s.IsType(test.expectedSubject, subjectNode)
			s.IsType(test.expectedNode, node)
		})
	}
}
