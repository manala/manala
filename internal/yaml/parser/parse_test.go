package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/yaml/parser"

	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestEmpty() {
	node, err := parser.Parse(nil)

	s.Nil(node)

	errors.Equal(s.T(), &serrors.Assertion{
		Message: "empty yaml file",
	}, err)
}

func (s *Suite) TestInvalids() {
	tests := []struct {
		test     string
		expected *serrors.Assertion
	}{
		{
			test: "At",
			expected: &serrors.Assertion{
				Message: "'@' is a reserved character",
				Arguments: []any{
					"line", 1,
					"column", 1,
				},
				Details: `
					>  1 | @
					       ^
				`,
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

	errors.Equal(s.T(), &serrors.Assertion{
		Message: "multiple documents yaml file",
		Arguments: []any{
			"line", 4,
			"column", 1,
		},
		Details: `
			   1 | ---
			   2 | foo
			   3 | ---
			>  4 | bar
			       ^
		`,
	}, err)
}

func (s *Suite) TestIrregularMapKeys() {
	tests := []struct {
		test     string
		expected *serrors.Assertion
	}{
		{
			test: "Integer",
			expected: &serrors.Assertion{
				Message: "irregular map key",
				Arguments: []any{
					"line", 1,
					"column", 2,
				},
				Details: `
					>  1 | 0: foo
					        ^
				`,
			},
		},
		{
			test: "IntegerAnchor",
			expected: &serrors.Assertion{
				Message: "irregular map key",
				Arguments: []any{
					"line", 2,
					"column", 4,
				},
				Details: `
					   1 | anchor: &anchor
					>  2 |   0: foo
					          ^
				`,
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
		expected *serrors.Assertion
	}{
		{
			test: "Inf",
			expected: &serrors.Assertion{
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
		{
			test: "Nan",
			expected: &serrors.Assertion{
				Message: "irregular type",
				Arguments: []any{
					"line", 1,
					"column", 6,
				},
				Details: `
					>  1 | foo: .nan
					            ^
				`,
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
	s.Require().IsType((*goYamlAst.MappingKeyNode)(nil), keyNode)
	keyNodeValue := keyNode.(*goYamlAst.MappingKeyNode).Value
	s.Require().IsType((*goYamlAst.StringNode)(nil), keyNodeValue)
	s.Equal("foo", keyNodeValue.(*goYamlAst.StringNode).Value)

	valueNode := node.Values[0].Value
	s.Require().IsType((*goYamlAst.StringNode)(nil), valueNode)
	s.Equal("bar", valueNode.(*goYamlAst.StringNode).Value)
}

func (s *Suite) TestIrregularMappingKey() {
	dir := filepath.FromSlash("testdata/TestIrregularMappingKey")
	content, _ := os.ReadFile(filepath.Join(dir, "node.yaml"))

	node, err := parser.Parse(content)

	s.Nil(node)

	errors.Equal(s.T(), &serrors.Assertion{
		Message: "irregular map key",
		Arguments: []any{
			"line", 1,
			"column", 1,
		},
		Details: `
			>  1 | ? 123: bar
			       ^
		`,
	}, err)
}

func (s *Suite) TestTags() {
	dir := filepath.FromSlash("testdata/TestTags")
	content, _ := os.ReadFile(filepath.Join(dir, "node.yaml"))

	node, err := parser.Parse(content)

	s.Require().NoError(err)

	s.Require().Len(node.Values, 1)

	s.Require().IsType((*goYamlAst.StringNode)(nil), node.Values[0].Value)
	s.Equal("bar", node.Values[0].Value.String())
}

func (s *Suite) TestUnknownAnchors() {
	dir := filepath.FromSlash("testdata/TestUnknownAnchors")
	content, _ := os.ReadFile(filepath.Join(dir, "node.yaml"))

	node, err := parser.Parse(content)

	s.Nil(node)

	errors.Equal(s.T(), &serrors.Assertion{
		Message: "cannot find anchor",
		Arguments: []any{
			"line", 1,
			"column", 7,
			"anchor", "anchor",
		},
		Details: `
			>  1 | foo: *anchor
			             ^
		`,
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
		s.Require().IsType((*goYamlAst.StringNode)(nil), anchorNode.Value)
		s.Equal("foo", anchorNode.Value.(*goYamlAst.StringNode).Value)

		aliasNode := node.Values[1]
		s.Require().IsType((*goYamlAst.StringNode)(nil), aliasNode.Value)
		s.Equal("foo", aliasNode.Value.(*goYamlAst.StringNode).Value)
	})
	s.Run("MergeKeys", func() {
		dir := filepath.FromSlash("testdata/TestAnchors")
		content, _ := os.ReadFile(filepath.Join(dir, "MergeKeys.yaml"))

		node, err := parser.Parse(content)

		s.Require().NoError(err)

		s.Require().Len(node.Values, 9)

		emptyAnchorNode := node.Values[0]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), emptyAnchorNode.Value)
		s.Empty(emptyAnchorNode.Value.(*goYamlAst.MappingNode).Values)

		mappingValueAnchorNode := node.Values[1]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingValueAnchorNode.Value)
		s.Require().Len(mappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values, 1)

		mappingAnchorNode := node.Values[2]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingAnchorNode.Value)
		s.Require().Len(mappingAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)

		mappingValueAliasEmptyAnchorNode := node.Values[3]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingValueAliasEmptyAnchorNode.Value)
		s.Empty(mappingValueAliasEmptyAnchorNode.Value.(*goYamlAst.MappingNode).Values)

		mappingValueAliasMappingValueAnchorNode := node.Values[4]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingValueAliasMappingValueAnchorNode.Value)
		s.Require().Len(mappingValueAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values, 1)

		mappingValueAliasMappingAnchorNode := node.Values[5]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingValueAliasMappingAnchorNode.Value)
		s.Require().Len(mappingValueAliasMappingAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)

		mappingAliasEmptyAnchorNode := node.Values[6]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingAliasEmptyAnchorNode.Value)
		s.Require().Len(mappingAliasEmptyAnchorNode.Value.(*goYamlAst.MappingNode).Values, 1)

		mappingAliasMappingValueAnchorNode := node.Values[7]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingAliasMappingValueAnchorNode.Value)
		s.Require().Len(mappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)

		mappingValueAliasMappingNode := node.Values[8]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingValueAliasMappingNode.Value)
		s.Require().Len(mappingValueAliasMappingNode.Value.(*goYamlAst.MappingNode).Values, 3)
	})
	s.Run("MergeKeysDuplicated", func() {
		dir := filepath.FromSlash("testdata/TestAnchors")
		content, _ := os.ReadFile(filepath.Join(dir, "MergeKeysDuplicated.yaml"))

		node, err := parser.Parse(content)

		s.Require().NoError(err)

		s.Require().Len(node.Values, 5)

		singleMappingAliasMappingValueAnchorNode := node.Values[2]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), singleMappingAliasMappingValueAnchorNode.Value)
		s.Require().Len(singleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values, 1)
		s.Equal("bar", singleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values[0].Value.(*goYamlAst.StringNode).Value)

		multipleMappingAliasMappingValueAnchorNode := node.Values[3]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), multipleMappingAliasMappingValueAnchorNode.Value)
		s.Require().Len(multipleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)
		s.Equal("bar", multipleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values[0].Value.(*goYamlAst.StringNode).Value)

		mappingAliasMappingAnchorNode := node.Values[4]
		s.Require().IsType((*goYamlAst.MappingNode)(nil), mappingAliasMappingAnchorNode.Value)
		s.Require().Len(mappingAliasMappingAnchorNode.Value.(*goYamlAst.MappingNode).Values, 3)
		s.Equal("bar", multipleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values[1].Value.(*goYamlAst.StringNode).Value)
	})
}
