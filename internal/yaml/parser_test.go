package yaml

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	"manala/internal/serrors"
	"path/filepath"
)

func (s *Suite) TestParserEmpty() {
	parser := NewParser()

	node, err := parser.ParseBytes(nil)

	s.Nil(node)

	serrors.Equal(s.Assert(), &serrors.Assert{
		Type:    serrors.Error{},
		Message: "empty yaml file",
	}, err)
}

func (s *Suite) TestParserMultipleDocuments() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestMultipleDocuments")

	parser := NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Nil(node)

	serrors.Equal(s.Assert(), &serrors.Assert{
		Type:    serrors.Error{},
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

func (s *Suite) TestParserMappingComments() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestMappingComments")

	parser := NewParser(WithComments())
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.NoError(err)

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

func (s *Suite) TestParserIrregularMapKeys() {
	tests := []struct {
		test     string
		expected *serrors.Assert
	}{
		{
			test: "Integer",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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

			parser := NewParser()
			node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			s.Nil(node)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestParserIrregularTypes() {
	tests := []struct {
		test     string
		expected *serrors.Assert
	}{
		{
			test: "Inf",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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

			parser := NewParser()
			node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

			s.Nil(node)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestParserMappingKey() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestMappingKey")

	parser := NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.NoError(err)

	s.IsType((*goYamlAst.MappingValueNode)(nil), node)

	keyNode := node.(*goYamlAst.MappingValueNode).Key
	s.IsType((*goYamlAst.StringNode)(nil), keyNode)
	s.Equal("foo", keyNode.(*goYamlAst.StringNode).Value)

	valueNode := node.(*goYamlAst.MappingValueNode).Value
	s.IsType((*goYamlAst.StringNode)(nil), valueNode)
	s.Equal("bar", valueNode.(*goYamlAst.StringNode).Value)
}

func (s *Suite) TestParserIrregularMappingKey() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestIrregularMappingKey")

	parser := NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Nil(node)

	serrors.Equal(s.Assert(), &serrors.Assert{
		Type:    serrors.Error{},
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

func (s *Suite) TestParserTags() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestTags")

	parser := NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.NoError(err)

	s.IsType((*goYamlAst.StringNode)(nil), node)
	s.Equal("foo", node.(*goYamlAst.StringNode).Value)
}

func (s *Suite) TestParserUnknownAnchors() {
	dir := filepath.FromSlash("testdata/ParserSuite/TestUnknownAnchors")

	parser := NewParser()
	node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

	s.Nil(node)

	serrors.Equal(s.Assert(), &serrors.Assert{
		Type:    serrors.Error{},
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

func (s *Suite) TestParserAnchors() {
	s.Run("Anchors", func() {
		dir := filepath.FromSlash("testdata/ParserSuite/TestAnchors/Anchors")

		parser := NewParser()
		node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

		s.NoError(err)

		anchorNode := node.(*goYamlAst.MappingNode).Values[0]
		s.IsType((*goYamlAst.StringNode)(nil), anchorNode.Value)
		s.Equal("foo", anchorNode.Value.(*goYamlAst.StringNode).Value)

		aliasNode := node.(*goYamlAst.MappingNode).Values[1]
		s.IsType((*goYamlAst.StringNode)(nil), aliasNode.Value)
		s.Equal("foo", aliasNode.Value.(*goYamlAst.StringNode).Value)
	})
	s.Run("MergeKeys", func() {
		dir := filepath.FromSlash("testdata/ParserSuite/TestAnchors/MergeKeys")

		parser := NewParser()
		node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

		s.NoError(err)

		emptyAnchorNode := node.(*goYamlAst.MappingNode).Values[0]
		s.IsType((*goYamlAst.MappingNode)(nil), emptyAnchorNode.Value)
		s.Len(emptyAnchorNode.Value.(*goYamlAst.MappingNode).Values, 0)

		mappingValueAnchorNode := node.(*goYamlAst.MappingNode).Values[1]
		s.IsType((*goYamlAst.MappingValueNode)(nil), mappingValueAnchorNode.Value)

		mappingAnchorNode := node.(*goYamlAst.MappingNode).Values[2]
		s.IsType((*goYamlAst.MappingNode)(nil), mappingAnchorNode.Value)
		s.Len(mappingAnchorNode.Value.(*goYamlAst.MappingNode).Values, 2)

		mappingValueAliasEmptyAnchorNode := node.(*goYamlAst.MappingNode).Values[3]
		s.IsType((*goYamlAst.MappingNode)(nil), mappingValueAliasEmptyAnchorNode.Value)
		s.Len(mappingValueAliasEmptyAnchorNode.Value.(*goYamlAst.MappingNode).Values, 0)

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

		parser := NewParser()
		node, err := parser.ParseFile(filepath.Join(dir, "node.yaml"))

		s.NoError(err)

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
