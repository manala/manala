package parser_test

import (
	"testing"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/expect"
	"github.com/manala/manala/internal/testing/heredoc"
	"github.com/manala/manala/internal/yaml/parser"

	"github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type ParseSuite struct{ suite.Suite }

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func (s *ParseSuite) TestAnchors() {
	s.Run("Anchors", func() {
		node, err := parser.Parse([]byte(heredoc.Doc(`
			anchor: &anchor foo
			alias: *anchor
		`)))

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
		node, err := parser.Parse([]byte(heredoc.Doc(`
			empty: &empty {}
			mapping_value: &mapping_value
			  foo: foo
			mapping: &mapping
			  foo: foo
			  bar: bar
			mapping_value_alias_empty:
			  <<: *empty
			mapping_value_alias_mapping_value:
			  <<: *mapping_value
			mapping_value_alias_mapping:
			  <<: *mapping
			mapping_alias_empty:
			  <<: *empty
			  baz: baz
			mapping_alias_mapping_value:
			  <<: *mapping_value
			  baz: baz
			mapping_alias_mapping:
			  <<: *mapping
			  baz: baz
		`)))

		s.Require().NoError(err)

		s.Require().Len(node.Values, 9)

		emptyNode := node.Values[0]
		s.Require().IsType((*ast.MappingNode)(nil), emptyNode.Value)
		s.Empty(emptyNode.Value.(*ast.MappingNode).Values)

		mappingValueNode := node.Values[1]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueNode.Value)
		s.Require().Len(mappingValueNode.Value.(*ast.MappingNode).Values, 1)

		mappingNode := node.Values[2]
		s.Require().IsType((*ast.MappingNode)(nil), mappingNode.Value)
		s.Require().Len(mappingNode.Value.(*ast.MappingNode).Values, 2)

		mappingValueAliasEmptyNode := node.Values[3]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueAliasEmptyNode.Value)
		s.Empty(mappingValueAliasEmptyNode.Value.(*ast.MappingNode).Values)

		mappingValueAliasMappingValueNode := node.Values[4]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueAliasMappingValueNode.Value)
		s.Require().Len(mappingValueAliasMappingValueNode.Value.(*ast.MappingNode).Values, 1)

		mappingValueAliasMappingNode := node.Values[5]
		s.Require().IsType((*ast.MappingNode)(nil), mappingValueAliasMappingNode.Value)
		s.Require().Len(mappingValueAliasMappingNode.Value.(*ast.MappingNode).Values, 2)

		mappingAliasEmptyNode := node.Values[6]
		s.Require().IsType((*ast.MappingNode)(nil), mappingAliasEmptyNode.Value)
		s.Require().Len(mappingAliasEmptyNode.Value.(*ast.MappingNode).Values, 1)

		mappingAliasMappingValueNode := node.Values[7]
		s.Require().IsType((*ast.MappingNode)(nil), mappingAliasMappingValueNode.Value)
		s.Require().Len(mappingAliasMappingValueNode.Value.(*ast.MappingNode).Values, 2)

		mappingAliasMappingNode := node.Values[8]
		s.Require().IsType((*ast.MappingNode)(nil), mappingAliasMappingNode.Value)
		s.Require().Len(mappingAliasMappingNode.Value.(*ast.MappingNode).Values, 3)
	})

	s.Run("MergeKeysDuplicated", func() {
		node, err := parser.Parse([]byte(heredoc.Doc(`
			mapping_value: &mapping_value
			  foo: foo
			mapping: &mapping
			  foo: foo
			  bar: bar
			single_mapping_alias_mapping_value:
			  <<: *mapping_value
			  foo: bar
			multiple_mapping_alias_mapping_value:
			  <<: *mapping_value
			  foo: bar
			  bar: bar
			mapping_alias_mapping:
			  <<: *mapping
			  foo: bar
			  baz: baz
		`)))

		s.Require().NoError(err)

		s.Require().Len(node.Values, 5)

		singleMappingAliasMappingValueNode := node.Values[2]
		s.Require().IsType((*ast.MappingNode)(nil), singleMappingAliasMappingValueNode.Value)
		s.Require().Len(singleMappingAliasMappingValueNode.Value.(*ast.MappingNode).Values, 1)
		s.Equal("bar", singleMappingAliasMappingValueNode.Value.(*ast.MappingNode).Values[0].Value.(*ast.StringNode).Value)

		multipleMappingAliasMappingValueNode := node.Values[3]
		s.Require().IsType((*ast.MappingNode)(nil), multipleMappingAliasMappingValueNode.Value)
		s.Require().Len(multipleMappingAliasMappingValueNode.Value.(*ast.MappingNode).Values, 2)
		s.Equal("bar", multipleMappingAliasMappingValueNode.Value.(*ast.MappingNode).Values[0].Value.(*ast.StringNode).Value)

		mappingAliasMappingNode := node.Values[4]
		s.Require().IsType((*ast.MappingNode)(nil), mappingAliasMappingNode.Value)
		s.Require().Len(mappingAliasMappingNode.Value.(*ast.MappingNode).Values, 3)
		s.Equal("bar", multipleMappingAliasMappingValueNode.Value.(*ast.MappingNode).Values[1].Value.(*ast.StringNode).Value)
	})
}

func (s *ParseSuite) TestTags() {
	node, err := parser.Parse([]byte(heredoc.Doc(`
		foo: !!str bar
	`)))

	s.Require().NoError(err)

	s.Require().Len(node.Values, 1)

	s.Require().IsType((*ast.StringNode)(nil), node.Values[0].Value)
	s.Equal("bar", node.Values[0].Value.String())
}

func (s *ParseSuite) TestMapKeyExplicit() {
	node, err := parser.Parse([]byte(heredoc.Doc(`
		? foo: bar
	`)))

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

func (s *ParseSuite) TestErrors() {
	tests := []struct {
		test     string
		data     string
		expected expect.ErrorExpectation
	}{
		{
			test: "Empty",
			data: "",
			expected: serrors.Expectation{
				Message: "empty yaml content",
			},
		},
		{
			test: "Spaces",
			data: " ",
			expected: serrors.Expectation{
				Message: "empty yaml content",
			},
		},
		{
			test: "Invalid",
			data: "@",
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("'@' is a reserved character"),
			},
		},
		{
			test: "Tab",
			data: heredoc.Doc(`
				foo: bar
					baz: qux
			`),
			expected: parsing.ErrorExpectation{
				Line:   2,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("found a tab character where an indentation space is expected "),
			},
		},
		{
			test: "MultipleDocuments",
			data: heredoc.Doc(`
				---
				foo
				---
				bar
			`),
			expected: parsing.ErrorExpectation{
				Line:   3,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("multiple documents yaml content"),
			},
		},
		{
			test: "NotMap",
			data: heredoc.Doc(`
				foo
			`),
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("yaml document must be a map"),
			},
		},
		{
			test: "IrregularTypeInf",
			data: heredoc.Doc(`
				foo: .inf
			`),
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 6,
				Err:    expect.ErrorMessageExpectation("irregular yaml type"),
			},
		},
		{
			test: "IrregularTypeNan",
			data: heredoc.Doc(`
				foo: .nan
			`),
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 6,
				Err:    expect.ErrorMessageExpectation("irregular yaml type"),
			},
		},
		{
			test: "IrregularMapKeyInteger",
			data: heredoc.Doc(`
				123: foo
			`),
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("irregular yaml map key"),
			},
		},
		{
			test: "IrregularMapKeyExplicit",
			data: heredoc.Doc(`
				? 123: bar
			`),
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("irregular yaml map key"),
			},
		},
		{
			test: "IrregularMapKeyInteger",
			data: heredoc.Doc(`
				anchor: &anchor
				  0: foo
			`),
			expected: parsing.ErrorExpectation{
				Line:   2,
				Column: 3,
				Err:    expect.ErrorMessageExpectation("irregular yaml map key"),
			},
		},
		{
			test: "UnknownAnchors",
			data: heredoc.Doc(`
				foo: *bar
			`),
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 7,
				Err:    expect.ErrorMessageExpectation("unknown \"bar\" yaml anchor"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := parser.Parse([]byte(test.data))

			s.Nil(node)
			expect.Error(s.T(), test.expected, err)
		})
	}
}
