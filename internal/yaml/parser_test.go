package yaml_test

import (
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/yaml"
	"path/filepath"
	"testing"

	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type ParserSuite struct{ suite.Suite }

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}

func (s *ParserSuite) TestEmpty() {
	parser := yaml.NewParser()

	node, err := parser.ParseBytes(nil)

	s.Nil(node)

	serrors.Equal(s.T(), &serrors.Assertion{
		Message: "empty yaml file",
	}, err)
}

func (s *ParserSuite) TestMultipleDocuments() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestMultipleDocuments")

	parser := yaml.NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Nil(node)

	serrors.Equal(s.T(), &serrors.Assertion{
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

func (s *ParserSuite) TestMappingComments() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestMappingComments")

	parser := yaml.NewParser(yaml.WithComments())
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Require().NoError(err)

	emptyNode := node.(*goYamlAst.MappingNode).Values[0]
	s.Equal("# Empty", emptyNode.GetComment().String())

	mappingValueNode := node.(*goYamlAst.MappingNode).Values[1]
	s.Equal("# Mapping Value", mappingValueNode.GetComment().String())
	s.Equal("# Mapping Value Foo", mappingValueNode.Value.GetComment().String())

	mappingNode := node.(*goYamlAst.MappingNode).Values[2]
	s.Equal("# Mapping", mappingNode.GetComment().String())
	s.Equal("# Mapping Foo", mappingNode.Value.(*goYamlAst.MappingNode).Values[0].GetComment().String())
	s.Equal("# Mapping Bar", mappingNode.Value.(*goYamlAst.MappingNode).Values[1].GetComment().String())
}

func (s *ParserSuite) TestIrregularMapKeys() {
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
			dir := filepath.FromSlash("testdata/ParserSuite/TestIrregularMapKeys/" + test.test)

			parser := yaml.NewParser()
			node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			s.Nil(node)

			serrors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *ParserSuite) TestIrregularTypes() {
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
			dir := filepath.FromSlash("testdata/ParserSuite/TestIrregularTypes/" + test.test)

			parser := yaml.NewParser()
			node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			s.Nil(node)

			serrors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *ParserSuite) TestMappingKey() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestMappingKey")

	parser := yaml.NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Require().NoError(err)

	s.IsType((*goYamlAst.MappingValueNode)(nil), node)

	keyNode := node.(*goYamlAst.MappingValueNode).Key
	s.IsType((*goYamlAst.StringNode)(nil), keyNode)
	s.Equal("foo", keyNode.(*goYamlAst.StringNode).Value)

	valueNode := node.(*goYamlAst.MappingValueNode).Value
	s.IsType((*goYamlAst.StringNode)(nil), valueNode)
	s.Equal("bar", valueNode.(*goYamlAst.StringNode).Value)
}

func (s *ParserSuite) TestIrregularMappingKey() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestIrregularMappingKey")

	parser := yaml.NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Nil(node)

	serrors.Equal(s.T(), &serrors.Assertion{
		Message: "irregular map key",
		Arguments: []any{
			"line", 1,
			"column", 6,
		},
		Details: `
			>  1 | ? 123: bar
			            ^
		`,
	}, err)
}

func (s *ParserSuite) TestLiteralString() {
	tests := []struct {
		test     string
		expected string
	}{
		{
			test:     "SingleTrailingNoStrip",
			expected: "bar\n",
		},
		{
			test:     "SingleTrailingStrip",
			expected: "bar",
		},
		{
			test:     "MultipleTrailingsNoStrip",
			expected: "bar\n",
		},
		{
			test:     "MultipleTrailingsStrip",
			expected: "bar",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			dir := filepath.FromSlash("testdata/ParserSuite/TestLiteralString/" + test.test)

			parser := yaml.NewParser()
			node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			s.IsType((*goYamlAst.MappingValueNode)(nil), node)
			value := node.(*goYamlAst.MappingValueNode).
				Value.(*goYamlAst.LiteralNode).
				Value.Value

			s.Require().NoError(err)
			s.Equal(test.expected, value)
		})
	}
}

func (s *ParserSuite) TestTags() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestTags")

	parser := yaml.NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Require().NoError(err)

	s.IsType((*goYamlAst.StringNode)(nil), node)
	s.Equal("foo", node.(*goYamlAst.StringNode).Value)
}

func (s *ParserSuite) TestUnknownAnchors() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestUnknownAnchors")

	parser := yaml.NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Nil(node)

	serrors.Equal(s.T(), &serrors.Assertion{
		Message: "cannot find anchor",
		Arguments: []any{
			"line", 1,
			"column", 2,
			"anchor", "anchor",
		},
		Details: `
			>  1 | *anchor
			        ^
		`,
	}, err)
}

func (s *ParserSuite) TestAnchors() {
	s.Run("Anchors", func() {
		dir := filepath.FromSlash("testdata/ParserSuite/TestAnchors/Anchors")

		parser := yaml.NewParser()
		node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

		s.Require().NoError(err)

		anchorNode := node.(*goYamlAst.MappingNode).Values[0]
		s.IsType((*goYamlAst.StringNode)(nil), anchorNode.Value)
		s.Equal("foo", anchorNode.Value.(*goYamlAst.StringNode).Value)

		aliasNode := node.(*goYamlAst.MappingNode).Values[1]
		s.IsType((*goYamlAst.StringNode)(nil), aliasNode.Value)
		s.Equal("foo", aliasNode.Value.(*goYamlAst.StringNode).Value)
	})
	s.Run("MergeKeys", func() {
		dir := filepath.FromSlash("testdata/ParserSuite/TestAnchors/MergeKeys")

		parser := yaml.NewParser()
		node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

		s.Require().NoError(err)

		emptyAnchorNode := node.(*goYamlAst.MappingNode).Values[0]
		s.IsType((*goYamlAst.MappingNode)(nil), emptyAnchorNode.Value)
		s.Empty(emptyAnchorNode.Value.(*goYamlAst.MappingNode).Values)

		mappingValueAnchorNode := node.(*goYamlAst.MappingNode).Values[1]
		s.IsType((*goYamlAst.MappingValueNode)(nil), mappingValueAnchorNode.Value)

		mappingAnchorNode := node.(*goYamlAst.MappingNode).Values[2]
		s.IsType((*goYamlAst.MappingNode)(nil), mappingAnchorNode.Value)
		s.Len(mappingAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)

		mappingValueAliasEmptyAnchorNode := node.(*goYamlAst.MappingNode).Values[3]
		s.IsType((*goYamlAst.MappingNode)(nil), mappingValueAliasEmptyAnchorNode.Value)
		s.Empty(mappingValueAliasEmptyAnchorNode.Value.(*goYamlAst.MappingNode).Values)

		mappingValueAliasMappingValueAnchorNode := node.(*goYamlAst.MappingNode).Values[4]
		s.IsType((*goYamlAst.MappingValueNode)(nil), mappingValueAliasMappingValueAnchorNode.Value)

		mappingValueAliasMappingAnchorNode := node.(*goYamlAst.MappingNode).Values[5]
		s.IsType((*goYamlAst.MappingNode)(nil), mappingValueAliasMappingAnchorNode.Value)
		s.Len(mappingValueAliasMappingAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)

		mappingAliasEmptyAnchorNode := node.(*goYamlAst.MappingNode).Values[6]
		s.IsType((*goYamlAst.MappingValueNode)(nil), mappingAliasEmptyAnchorNode.Value)

		mappingAliasMappingValueAnchorNode := node.(*goYamlAst.MappingNode).Values[7]
		s.IsType((*goYamlAst.MappingNode)(nil), mappingAliasMappingValueAnchorNode.Value)
		s.Len(mappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)

		mappingValueAliasMappingNode := node.(*goYamlAst.MappingNode).Values[8]
		s.IsType((*goYamlAst.MappingNode)(nil), mappingValueAliasMappingNode.Value)
		s.Len(mappingValueAliasMappingNode.Value.(*goYamlAst.MappingNode).Values, 3)
	})
	s.Run("MergeKeysDuplicated", func() {
		dir := filepath.FromSlash("testdata/ParserSuite/TestAnchors/MergeKeysDuplicated")

		parser := yaml.NewParser()
		node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

		s.Require().NoError(err)

		singleMappingAliasMappingValueAnchorNode := node.(*goYamlAst.MappingNode).Values[2]
		s.IsType((*goYamlAst.MappingValueNode)(nil), singleMappingAliasMappingValueAnchorNode.Value)
		s.Equal("bar", singleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingValueNode).Value.(*goYamlAst.StringNode).Value)

		multipleMappingAliasMappingValueAnchorNode := node.(*goYamlAst.MappingNode).Values[3]
		s.IsType((*goYamlAst.MappingNode)(nil), multipleMappingAliasMappingValueAnchorNode.Value)
		s.Len(multipleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)
		s.Equal("bar", multipleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values[0].Value.(*goYamlAst.StringNode).Value)

		mappingAliasMappingAnchorNode := node.(*goYamlAst.MappingNode).Values[4]
		s.IsType((*goYamlAst.MappingNode)(nil), mappingAliasMappingAnchorNode.Value)
		s.Len(mappingAliasMappingAnchorNode.Value.(*goYamlAst.MappingNode).Values, 3)
		s.Equal("bar", multipleMappingAliasMappingValueAnchorNode.Value.(*goYamlAst.MappingNode).Values[1].Value.(*goYamlAst.StringNode).Value)
	})
}
