package yaml

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
	"manala/internal/errors/serrors"
	"manala/internal/testing/heredoc"
	"path/filepath"
	"testing"
)

type ExtractorSuite struct{ suite.Suite }

func TestExtractorSuite(t *testing.T) {
	suite.Run(t, new(ExtractorSuite))
}

func (s *ExtractorSuite) TestExtractRootMapErrors() {
	tests := []struct {
		test     string
		expected *serrors.Assert
	}{
		{
			test: "Empty",
			expected: &serrors.Assert{
				Type:    &NodeError{},
				Message: "root must be a map",
			},
		},
		{
			test: "NonMap",
			expected: &serrors.Assert{
				Type:    &NodeError{},
				Message: "root must be a map",
				Arguments: []any{
					"line", 1,
					"column", 1,
				},
				Details: heredoc.Doc(`
					>  1 | foo
					       ^
				`),
			},
		},
		{
			test: "SubjectNotFoundSingle",
			expected: &serrors.Assert{
				Type:    &serrors.Error{},
				Message: "unable to find map",
				Arguments: []any{
					"key", "subject",
				},
			},
		},
		{
			test: "SubjectNotFoundMultiple",
			expected: &serrors.Assert{
				Type:    &serrors.Error{},
				Message: "unable to find map",
				Arguments: []any{
					"key", "subject",
				},
			},
		},
		{
			test: "SubjectNonMapSingle",
			expected: &serrors.Assert{
				Type:    &NodeError{},
				Message: "key is not a map",
				Arguments: []any{
					"line", 1,
					"column", 10,
					"key", "subject",
				},
				Details: heredoc.Doc(`
					>  1 | subject: 123
					                ^
				`),
			},
		},
		{
			test: "SubjectNonMapMultiple",
			expected: &serrors.Assert{
				Type:    &NodeError{},
				Message: "key is not a map",
				Arguments: []any{
					"line", 1,
					"column", 10,
					"key", "subject",
				},
				Details: heredoc.Doc(`
					>  1 | subject: 123
					                ^
					   2 | foo: foo
				`),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			dir := filepath.FromSlash("testdata/ExtractorSuite/TestExtractRootMapErrors/" + test.test)

			parser := NewParser()
			node, _ := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			extractor := NewExtractor(&node)
			subjectNode, err := extractor.ExtractRootMap("subject")

			s.Nil(subjectNode)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *ExtractorSuite) TestExtractRootMap() {
	tests := []struct {
		test            string
		expectedSubject interface{}
		expectedNode    interface{}
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
			dir := filepath.FromSlash("testdata/ExtractorSuite/TestExtractRootMap/" + test.test)

			parser := NewParser()
			node, _ := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			extractor := NewExtractor(&node)
			subjectNode, err := extractor.ExtractRootMap("subject")

			s.NoError(err)

			s.IsType(test.expectedSubject, subjectNode)
			s.IsType(test.expectedNode, node)
		})
	}
}
