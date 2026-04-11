package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/yaml/parser"

	"github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestEmpty() {
	node, err := parser.Parse(nil)

	s.Nil(node)

	errors.Equal(s.T(), &parsing.ErrorAssertion{
		Err: &serrors.Assertion{
			Message: "empty yaml content",
		},
	}, err)
}

func (s *Suite) TestInvalids() {
	tests := []struct {
		test     string
		expected errors.Assertion
	}{
		{
			test: "At",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "'@' is a reserved character",
				},
			},
		},
		{
			test: "Tab",
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "found a tab character where an indentation space is expected ",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			dir := filepath.FromSlash("testdata/TestInvalids")
			content, _ := os.ReadFile(filepath.Join(dir, test.test+".yaml"))

			node, err := parser.Parse(content)

			s.Nil(node)
			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *Suite) TestMultipleDocuments() {
	dir := filepath.FromSlash("testdata/TestMultipleDocuments")
	content, _ := os.ReadFile(filepath.Join(dir, "node.yaml"))

	node, err := parser.Parse(content)

	s.Nil(node)

	errors.Equal(s.T(), &parsing.ErrorAssertion{
		Line:   2,
		Column: 1,
		Err: &serrors.Assertion{
			Message: "multiple documents yaml content",
		},
	}, err)
}

func (s *Suite) TestIrregularMapKeys() {
	tests := []struct {
		test     string
		expected errors.Assertion
	}{
		{
			test: "Integer",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 2,
				Err: &serrors.Assertion{
					Message: "irregular map key",
				},
			},
		},
		{
			test: "IntegerAnchor",
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 4,
				Err: &serrors.Assertion{
					Message: "irregular map key",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			dir := filepath.FromSlash("testdata/TestIrregularMapKeys")
			content, _ := os.ReadFile(filepath.Join(dir, test.test+".yaml"))

			node, err := parser.Parse(content)

			s.Nil(node)
			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *Suite) TestIrregularTypes() {
	tests := []struct {
		test     string
		expected errors.Assertion
	}{
		{
			test: "Inf",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 6,
				Err: &serrors.Assertion{
					Message: "irregular type",
				},
			},
		},
		{
			test: "Nan",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 6,
				Err: &serrors.Assertion{
					Message: "irregular type",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			dir := filepath.FromSlash("testdata/TestIrregularTypes")
			content, _ := os.ReadFile(filepath.Join(dir, test.test+".yaml"))

			node, err := parser.Parse(content)

			s.Nil(node)
			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *Suite) TestMappingKey() {
	dir := filepath.FromSlash("testdata/TestMappingKey")
	content, _ := os.ReadFile(filepath.Join(dir, "node.yaml"))

	node, err := parser.Parse(content)

	s.Require().NoError(err)

	s.Require().Len(node.Values, 1)

	keyNode := node.Values[0].Key
	s.Require().IsType((*ast.MappingKeyNode)(nil), keyNode)
	keyNodeValue := keyNode.(*ast.MappingKeyNode).Value
	s.Require().IsType((*ast.StringNode)(nil), keyNodeValue)
	s.Equal("foo", keyNodeValue.(*ast.StringNode).Value)

	valueNode := node.Values[0].Value
	s.Require().IsType((*ast.StringNode)(nil), valueNode)
	s.Equal("bar", valueNode.(*ast.StringNode).Value)
}

func (s *Suite) TestIrregularMappingKey() {
	dir := filepath.FromSlash("testdata/TestIrregularMappingKey")
	content, _ := os.ReadFile(filepath.Join(dir, "node.yaml"))

	node, err := parser.Parse(content)

	s.Nil(node)

	errors.Equal(s.T(), &parsing.ErrorAssertion{
		Line:   1,
		Column: 1,
		Err: &serrors.Assertion{
			Message: "irregular map key",
		},
	}, err)
}

func (s *Suite) TestTags() {
	dir := filepath.FromSlash("testdata/TestTags")
	content, _ := os.ReadFile(filepath.Join(dir, "node.yaml"))

	node, err := parser.Parse(content)

	s.Require().NoError(err)

	s.Require().Len(node.Values, 1)

	s.Require().IsType((*ast.StringNode)(nil), node.Values[0].Value)
	s.Equal("bar", node.Values[0].Value.String())
}

func (s *Suite) TestUnknownAnchors() {
	dir := filepath.FromSlash("testdata/TestUnknownAnchors")
	content, _ := os.ReadFile(filepath.Join(dir, "node.yaml"))

	node, err := parser.Parse(content)

	s.Nil(node)

	errors.Equal(s.T(), &parsing.ErrorAssertion{
		Line:   1,
		Column: 7,
		Err: &serrors.Assertion{
			Message: "cannot find anchor",
			Arguments: []any{
				"anchor", "anchor",
			},
		},
	}, err)
}

func (s *Suite) TestAnchors() {
	s.Run("Anchors", func() {
		dir := filepath.FromSlash("testdata/TestAnchors")
		content, _ := os.ReadFile(filepath.Join(dir, "Anchors.yaml"))

		node, err := parser.Parse(content)

		s.Require().NoError(err)

		s.Require().Len(node.Values, 2)

		anchorNode := node.Values[0]
		s.Require().IsType((*ast.StringNode)(nil), anchorNode.Value)
		s.Equal("foo", anchorNode.Value.(*ast.StringNode).Value)

		aliasNode := node.Values[1]
		s.Require().IsType((*ast.StringNode)(nil), aliasNode.Value)
		s.Equal("foo", aliasNode.Value.(*ast.StringNode).Value)
	})
	s.Run("MergeKeys", func() {
		dir := filepath.FromSlash("testdata/TestAnchors")
		content, _ := os.ReadFile(filepath.Join(dir, "MergeKeys.yaml"))

		node, err := parser.Parse(content)

		s.Require().NoError(err)

		s.Require().Len(node.Values, 9)

		emptyAnchorNode := node.Values[0]
		s.Require().IsType((*ast.MappingNode)(nil), emptyAnchorNode.Value)
		s.Empty(emptyAnchorNode.Value.(*ast.MappingNode).Values)

		mappingValueAnchorNode := node.Values[1]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueAnchorNode.Value)
		s.Require().Len(mappingValueAnchorNode.Value.(*ast.MappingNode).Values, 1)

		mappingAnchorNode := node.Values[2]
		s.Require().IsType((*ast.MappingNode)(nil), mappingAnchorNode.Value)
		s.Require().Len(mappingAnchorNode.Value.(*ast.MappingNode).Values, 2)

		mappingValueAliasEmptyAnchorNode := node.Values[3]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueAliasEmptyAnchorNode.Value)
		s.Empty(mappingValueAliasEmptyAnchorNode.Value.(*ast.MappingNode).Values)

		mappingValueAliasMappingValueAnchorNode := node.Values[4]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueAliasMappingValueAnchorNode.Value)
		s.Require().Len(mappingValueAliasMappingValueAnchorNode.Value.(*ast.MappingNode).Values, 1)

		mappingValueAliasMappingAnchorNode := node.Values[5]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueAliasMappingAnchorNode.Value)
		s.Require().Len(mappingValueAliasMappingAnchorNode.Value.(*ast.MappingNode).Values, 2)

		mappingAliasEmptyAnchorNode := node.Values[6]
		s.Require().IsType((*ast.MappingNode)(nil), mappingAliasEmptyAnchorNode.Value)
		s.Require().Len(mappingAliasEmptyAnchorNode.Value.(*ast.MappingNode).Values, 1)

		mappingAliasMappingValueAnchorNode := node.Values[7]
		s.Require().IsType((*ast.MappingNode)(nil), mappingAliasMappingValueAnchorNode.Value)
		s.Require().Len(mappingAliasMappingValueAnchorNode.Value.(*ast.MappingNode).Values, 2)

		mappingValueAliasMappingNode := node.Values[8]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueAliasMappingNode.Value)
		s.Require().Len(mappingValueAliasMappingNode.Value.(*ast.MappingNode).Values, 3)
	})
	s.Run("MergeKeysDuplicated", func() {
		dir := filepath.FromSlash("testdata/TestAnchors")
		content, _ := os.ReadFile(filepath.Join(dir, "MergeKeysDuplicated.yaml"))

		node, err := parser.Parse(content)

		s.Require().NoError(err)

		s.Require().Len(node.Values, 5)

		singleMappingAliasMappingValueAnchorNode := node.Values[2]
		s.Require().IsType((*ast.MappingNode)(nil), singleMappingAliasMappingValueAnchorNode.Value)
		s.Require().Len(singleMappingAliasMappingValueAnchorNode.Value.(*ast.MappingNode).Values, 1)
		s.Equal("bar", singleMappingAliasMappingValueAnchorNode.Value.(*ast.MappingNode).Values[0].Value.(*ast.StringNode).Value)

		multipleMappingAliasMappingValueAnchorNode := node.Values[3]
		s.Require().IsType((*ast.MappingNode)(nil), multipleMappingAliasMappingValueAnchorNode.Value)
		s.Require().Len(multipleMappingAliasMappingValueAnchorNode.Value.(*ast.MappingNode).Values, 2)
		s.Equal("bar", multipleMappingAliasMappingValueAnchorNode.Value.(*ast.MappingNode).Values[0].Value.(*ast.StringNode).Value)

		mappingAliasMappingAnchorNode := node.Values[4]
		s.Require().IsType((*ast.MappingNode)(nil), mappingAliasMappingAnchorNode.Value)
		s.Require().Len(mappingAliasMappingAnchorNode.Value.(*ast.MappingNode).Values, 3)
		s.Equal("bar", multipleMappingAliasMappingValueAnchorNode.Value.(*ast.MappingNode).Values[1].Value.(*ast.StringNode).Value)
	})
}
