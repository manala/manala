package parser_test

import (
	"testing"

	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/testing/heredoc"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlparser "github.com/manala/manala/internal/yaml/parser"

	"github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type ParseSuite struct{ suite.Suite }

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func (s *ParseSuite) TestAnchors() {
	s.Run("Anchors", func() {
		node, err := yamlparser.Parse([]byte(heredoc.Doc(`
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
		node, err := yamlparser.Parse([]byte(heredoc.Doc(`
			empty: &empty {}
			single: &single
			  foo: foo
			multiple: &multiple
			  bar: bar
			  baz: baz

			mapping_empty_alias_empty:
			  <<: *empty
			mapping_empty_alias_single:
			  <<: *single
			mapping_empty_alias_multiple:
			  <<: *multiple
			mapping_empty_aliases_empty_single_multiple:
			  <<: [*empty, *single, *multiple]

			mapping_single_alias_empty:
			  <<: *empty
			  qux: qux
			mapping_single_alias_single:
			  <<: *single
			  qux: qux
			mapping_single_alias_multiple:
			  <<: *multiple
			  qux: qux
			mapping_single_aliases_empty_single_multiple:
			  <<: [*empty, *single, *multiple]
			  qux: qux
		`)))

		s.Require().NoError(err)

		s.Require().Len(node.Values, 11)

		empty := node.Values[0]
		s.Require().IsType((*ast.MappingNode)(nil), empty.Value)
		s.Empty(empty.Value.(*ast.MappingNode).Values)

		single := node.Values[1]
		s.Require().IsType((*ast.MappingNode)(nil), single.Value)
		s.Require().Len(single.Value.(*ast.MappingNode).Values, 1)

		multiple := node.Values[2]
		s.Require().IsType((*ast.MappingNode)(nil), multiple.Value)
		s.Require().Len(multiple.Value.(*ast.MappingNode).Values, 2)

		mappingEmptyAliasEmpty := node.Values[3]
		s.Require().IsType((*ast.MappingNode)(nil), mappingEmptyAliasEmpty.Value)
		s.Empty(mappingEmptyAliasEmpty.Value.(*ast.MappingNode).Values)

		mappingEmptyAliasSingle := node.Values[4]
		s.Require().IsType((*ast.MappingNode)(nil), mappingEmptyAliasSingle.Value)
		s.Require().Len(mappingEmptyAliasSingle.Value.(*ast.MappingNode).Values, 1)

		mappingEmptyAliasMultiple := node.Values[5]
		s.Require().IsType((*ast.MappingNode)(nil), mappingEmptyAliasMultiple.Value)
		s.Require().Len(mappingEmptyAliasMultiple.Value.(*ast.MappingNode).Values, 2)

		mappingEmptyAliasesEmptySingleMultiple := node.Values[6]
		s.Require().IsType((*ast.MappingNode)(nil), mappingEmptyAliasesEmptySingleMultiple.Value)
		s.Require().Len(mappingEmptyAliasesEmptySingleMultiple.Value.(*ast.MappingNode).Values, 3)

		mappingSingleAliasEmpty := node.Values[7]
		s.Require().IsType((*ast.MappingNode)(nil), mappingSingleAliasEmpty.Value)
		s.Require().Len(mappingSingleAliasEmpty.Value.(*ast.MappingNode).Values, 1)

		mappingSingleAliasSingle := node.Values[8]
		s.Require().IsType((*ast.MappingNode)(nil), mappingSingleAliasSingle.Value)
		s.Require().Len(mappingSingleAliasSingle.Value.(*ast.MappingNode).Values, 2)

		mappingSingleAliasMultiple := node.Values[9]
		s.Require().IsType((*ast.MappingNode)(nil), mappingSingleAliasMultiple.Value)
		s.Require().Len(mappingSingleAliasMultiple.Value.(*ast.MappingNode).Values, 3)

		mappingSingleAliasesEmptySingleMultiple := node.Values[10]
		s.Require().IsType((*ast.MappingNode)(nil), mappingSingleAliasesEmptySingleMultiple.Value)
		s.Require().Len(mappingSingleAliasesEmptySingleMultiple.Value.(*ast.MappingNode).Values, 4)
	})

	s.Run("MergeKeysDuplicated", func() {
		node, err := yamlparser.Parse([]byte(heredoc.Doc(`
			single: &single
			  foo: foo
			multiple: &multiple
			  foo: foo
			  bar: bar

			mapping_single_alias_single:
			  <<: *single
			  foo: bar
			mapping_single_aliases_single_multiple:
			  <<: [*single, *multiple]
			  foo: bar

			mapping_multiple_alias_single:
			  <<: *single
			  foo: bar
			  bar: bar
			mapping_multiple_alias_multiple:
			  <<: *multiple
			  foo: bar
			  baz: baz
			mapping_multiple_aliases_single_multiple:
			  <<: [*single, *multiple]
			  foo: bar
			  baz: baz
		`)))

		s.Require().NoError(err)

		s.Require().Len(node.Values, 7)

		mappingSingleAliasSingle := node.Values[2]
		s.Require().IsType((*ast.MappingNode)(nil), mappingSingleAliasSingle.Value)
		s.Require().Len(mappingSingleAliasSingle.Value.(*ast.MappingNode).Values, 1)
		s.Equal("bar", mappingSingleAliasSingle.Value.(*ast.MappingNode).Values[0].Value.(*ast.StringNode).Value)

		mappingSingleAliasesSingleMultiple := node.Values[3]
		s.Require().IsType((*ast.MappingNode)(nil), mappingSingleAliasesSingleMultiple.Value)
		s.Require().Len(mappingSingleAliasesSingleMultiple.Value.(*ast.MappingNode).Values, 2)
		s.Equal("bar", mappingSingleAliasesSingleMultiple.Value.(*ast.MappingNode).Values[0].Value.(*ast.StringNode).Value)

		mappingMultipleAliasSingle := node.Values[4]
		s.Require().IsType((*ast.MappingNode)(nil), mappingMultipleAliasSingle.Value)
		s.Require().Len(mappingMultipleAliasSingle.Value.(*ast.MappingNode).Values, 2)
		s.Equal("bar", mappingMultipleAliasSingle.Value.(*ast.MappingNode).Values[0].Value.(*ast.StringNode).Value)

		mappingMultipleAliasMultiple := node.Values[5]
		s.Require().IsType((*ast.MappingNode)(nil), mappingMultipleAliasMultiple.Value)
		s.Require().Len(mappingMultipleAliasMultiple.Value.(*ast.MappingNode).Values, 3)
		s.Equal("bar", mappingMultipleAliasMultiple.Value.(*ast.MappingNode).Values[1].Value.(*ast.StringNode).Value)

		mappingMultipleAliasesSingleMultiple := node.Values[6]
		s.Require().IsType((*ast.MappingNode)(nil), mappingMultipleAliasesSingleMultiple.Value)
		s.Require().Len(mappingMultipleAliasesSingleMultiple.Value.(*ast.MappingNode).Values, 3)
		s.Equal("bar", mappingMultipleAliasesSingleMultiple.Value.(*ast.MappingNode).Values[0].Value.(*ast.StringNode).Value)
	})
}

func (s *ParseSuite) TestAnchorsCycle() {
	s.Run("Anchor", func() {
		// Anchor cycles must surface as an explicit error rather than a stack overflow.
		_, err := yamlparser.Parse([]byte(heredoc.Doc(`
			cyclic: &cyclic
			  self: *cyclic
		`)))

		expectation.ExpectError(s.T(), yamlerrors.Expectation{
			Position: [2]int{2, 9},
			Err:      expectation.ErrorMessage("cycle through yaml anchor \"cyclic\""),
		}, err)
	})

	s.Run("MergeKey", func() {
		// Merge-key cycles must surface as an explicit error rather than a stack overflow.
		_, err := yamlparser.Parse([]byte(heredoc.Doc(`
			cyclic: &cyclic
			  <<: *cyclic
		`)))

		expectation.ExpectError(s.T(), yamlerrors.Expectation{
			Position: [2]int{2, 7},
			Err:      expectation.ErrorMessage("cycle through yaml anchor \"cyclic\""),
		}, err)
	})
}

func (s *ParseSuite) TestTags() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		foo: !!str bar
	`)))

	s.Require().NoError(err)

	s.Require().Len(node.Values, 1)

	value := node.Values[0].Value
	s.Require().IsType((*ast.StringNode)(nil), value)
	s.Equal("bar", value.String())
}

func (s *ParseSuite) TestMapKeyExplicit() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		? foo: bar
	`)))

	s.Require().NoError(err)

	s.Require().Len(node.Values, 1)

	key := node.Values[0].Key
	s.Require().IsType((*ast.MappingKeyNode)(nil), key)

	keyValue := key.(*ast.MappingKeyNode).Value
	s.Require().IsType((*ast.StringNode)(nil), keyValue)
	s.Equal("foo", keyValue.(*ast.StringNode).Value)

	value := node.Values[0].Value
	s.Require().IsType((*ast.StringNode)(nil), value)
	s.Equal("bar", value.(*ast.StringNode).Value)
}

func (s *ParseSuite) TestComments() {
	s.Run("PreservedThroughMerge", func() {
		// Merge without override: all comments must survive on merged keys.
		node, err := yamlparser.Parse([]byte(heredoc.Doc(`
			foo: &foo
			  # foo_foo annotation
			  foo_foo: 1
			  # foo_bar annotation
			  foo_bar: 2
			bar:
			  <<: *foo
		`)))
		s.Require().NoError(err)

		bar := node.Values[1].Value.(*ast.MappingNode)
		s.Require().Len(bar.Values, 2)

		s.Require().NotNil(bar.Values[0].GetComment())
		s.Equal("# foo_foo annotation", bar.Values[0].GetComment().String())

		s.Require().NotNil(bar.Values[1].GetComment())
		s.Equal("# foo_bar annotation", bar.Values[1].GetComment().String())
	})

	s.Run("InheritedOnOverrideWithoutChildComment", func() {
		// Override with no comment of its own should inherit the parent annotation.
		node, err := yamlparser.Parse([]byte(heredoc.Doc(`
				foo: &foo
				  # foo_foo annotation
				  foo_foo: 1
				bar:
				  <<: *foo
				  bar_foo: 2
			`)))
		s.Require().NoError(err)

		bar := node.Values[1].Value.(*ast.MappingNode)
		s.Require().Len(bar.Values, 2)

		s.Require().NotNil(bar.Values[0].GetComment())
		s.Equal("# foo_foo annotation", bar.Values[0].GetComment().String())
	})

	s.Run("OverriddenOnOverrideWithChildComment", func() {
		// Override that carries its own comment wins over the parent annotation.
		node, err := yamlparser.Parse([]byte(heredoc.Doc(`
				foo: &foo
				  # foo_foo annotation
				  foo_foo: 1
				bar:
				  <<: *foo
				  # bar_foo annotation
				  bar_foo: 2
			`)))
		s.Require().NoError(err)

		bar := node.Values[1].Value.(*ast.MappingNode)
		s.Require().Len(bar.Values, 2)

		s.Require().NotNil(bar.Values[0].GetComment())
		s.Equal("# foo_foo annotation", bar.Values[0].GetComment().String())

		s.Require().NotNil(bar.Values[1].GetComment())
		s.Equal("# bar_foo annotation", bar.Values[1].GetComment().String())
	})
}

func (s *ParseSuite) TestErrors() {
	tests := []struct {
		test     string
		data     string
		expected expectation.ErrorExpectation
	}{
		{
			test: "Empty",
			data: "",
			expected: yamlerrors.Expectation{
				Position: [2]int{0, 0},
				Err:      expectation.ErrorMessage("empty yaml content"),
			},
		},
		{
			test: "Spaces",
			data: " ",
			expected: yamlerrors.Expectation{
				Position: [2]int{0, 0},
				Err:      expectation.ErrorMessage("empty yaml content"),
			},
		},
		{
			test: "Invalid",
			data: "@",
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err:      expectation.ErrorMessage("'@' is a reserved character"),
			},
		},
		{
			test: "Tab",
			data: heredoc.Doc(`
				foo: bar
					baz: qux
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{2, 1},
				Err:      expectation.ErrorMessage("found a tab character where an indentation space is expected "),
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
			expected: yamlerrors.Expectation{
				Position: [2]int{3, 1},
				Err:      expectation.ErrorMessage("multiple documents yaml content"),
			},
		},
		{
			test: "NotMap",
			data: heredoc.Doc(`
				foo
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err:      expectation.ErrorMessage("yaml document must be a map"),
			},
		},
		{
			test: "IrregularTypeInf",
			data: heredoc.Doc(`
				foo: .inf
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 6},
				Err:      expectation.ErrorMessage("irregular yaml type"),
			},
		},
		{
			test: "IrregularTypeNan",
			data: heredoc.Doc(`
				foo: .nan
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 6},
				Err:      expectation.ErrorMessage("irregular yaml type"),
			},
		},
		{
			test: "IrregularMapKeyInteger",
			data: heredoc.Doc(`
				123: foo
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err:      expectation.ErrorMessage("irregular yaml map key"),
			},
		},
		{
			test: "IrregularMapKeyExplicit",
			data: heredoc.Doc(`
				? 123: bar
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err:      expectation.ErrorMessage("irregular yaml map key"),
			},
		},
		{
			test: "IrregularMapKeyInteger",
			data: heredoc.Doc(`
				anchor: &anchor
				  0: foo
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{2, 3},
				Err:      expectation.ErrorMessage("irregular yaml map key"),
			},
		},
		{
			test: "UnknownAnchors",
			data: heredoc.Doc(`
				foo: *bar
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 6},
				Err:      expectation.ErrorMessage("unknown \"bar\" yaml anchor"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := yamlparser.Parse([]byte(test.data))

			s.Nil(node)
			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
