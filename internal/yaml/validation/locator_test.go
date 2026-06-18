package validation_test

import (
	"testing"

	"github.com/manala/manala/internal/testing/heredoc"
	yamlvalidation "github.com/manala/manala/internal/yaml/validation"

	"github.com/goccy/go-yaml/parser"
	"github.com/stretchr/testify/suite"
)

type LocatorSuite struct{ suite.Suite }

func TestLocatorSuite(t *testing.T) {
	suite.Run(t, new(LocatorSuite))
}

func (s *LocatorSuite) TestAt() {
	tests := []struct {
		test     string
		node     string
		location string
		value    [2]int
		property [2]int
	}{
		{
			test: "Root",
			node: heredoc.Doc(`
				foo: bar
			`),
			location: "",
			value:    [2]int{0, 0},
			property: [2]int{1, 4},
		},
		{
			test: "NotFound",
			node: heredoc.Doc(`
				foo: bar
			`),
			location: "/baz",
			value:    [2]int{0, 0},
			property: [2]int{0, 0},
		},
		{
			test: "Mapping",
			node: heredoc.Doc(`
				foo: bar
			`),
			location: "/foo",
			value:    [2]int{1, 6},
			property: [2]int{1, 1},
		},
		{
			test: "MappingSecond",
			node: heredoc.Doc(`
				foo: bar
				baz: qux
			`),
			location: "/baz",
			value:    [2]int{2, 6},
			property: [2]int{2, 1},
		},
		{
			test: "MappingNested",
			node: heredoc.Doc(`
				foo:
				  bar: baz
			`),
			location: "/foo/bar",
			value:    [2]int{2, 8},
			property: [2]int{2, 3},
		},
		{
			test: "MappingMultiKey",
			node: heredoc.Doc(`
				foo:
				  bar: 1
				  baz: 2
			`),
			location: "/foo",
			value:    [2]int{1, 1},
			property: [2]int{1, 1},
		},
		{
			test: "MappingDeep",
			node: heredoc.Doc(`
				a:
				  b:
				    c: 1
			`),
			location: "/a/b/c",
			value:    [2]int{3, 8},
			property: [2]int{3, 5},
		},
		{
			test: "SequenceFirst",
			node: heredoc.Doc(`
				foo:
				  - 1
				  - 2
				  - 3
			`),
			location: "/foo/0",
			value:    [2]int{2, 5},
			property: [2]int{2, 5},
		},
		{
			test: "SequenceSecond",
			node: heredoc.Doc(`
				foo:
				  - 1
				  - 2
				  - 3
			`),
			location: "/foo/1",
			value:    [2]int{3, 5},
			property: [2]int{3, 5},
		},
		{
			test: "SequenceThird",
			node: heredoc.Doc(`
				foo:
				  - 1
				  - 2
				  - 3
			`),
			location: "/foo/2",
			value:    [2]int{4, 5},
			property: [2]int{4, 5},
		},
		{
			test: "SequenceOutOfBounds",
			node: heredoc.Doc(`
				foo:
				  - 1
				  - 2
				  - 3
			`),
			location: "/foo/5",
			value:    [2]int{0, 0},
			property: [2]int{0, 0},
		},
		{
			test: "RootSequenceFirst",
			node: heredoc.Doc(`
				- 1
				- 2
				- 3
			`),
			location: "/0",
			value:    [2]int{1, 3},
			property: [2]int{1, 3},
		},
		{
			test: "RootSequenceSecond",
			node: heredoc.Doc(`
				- 1
				- 2
				- 3
			`),
			location: "/1",
			value:    [2]int{2, 3},
			property: [2]int{2, 3},
		},
		{
			test: "SequenceOfMappings",
			node: heredoc.Doc(`
				- foo: 1
				- foo: 2
			`),
			location: "/1/foo",
			value:    [2]int{2, 8},
			property: [2]int{2, 3},
		},
		{
			test: "PointerEscapeTilde",
			node: heredoc.Doc(`
				a~b: 1
			`),
			location: "/a~0b",
			value:    [2]int{1, 6},
			property: [2]int{1, 1},
		},
		{
			test: "PointerEscapeSlash",
			node: heredoc.Doc(`
				a~b: 1
				a/b: 2
			`),
			location: "/a~1b",
			value:    [2]int{2, 6},
			property: [2]int{2, 1},
		},
		{
			test: "ValueBoolean",
			node: heredoc.Doc(`
				foo: true
			`),
			location: "/foo",
			value:    [2]int{1, 6},
			property: [2]int{1, 1},
		},
		{
			test: "ValueNull",
			node: heredoc.Doc(`
				foo: null
			`),
			location: "/foo",
			value:    [2]int{1, 6},
			property: [2]int{1, 1},
		},
		{
			test: "LeadingComment",
			node: heredoc.Doc(`
				# comment
				foo: bar
			`),
			location: "/foo",
			value:    [2]int{2, 6},
			property: [2]int{2, 1},
		},
		{
			test:     "FlowMapping",
			node:     `{foo: bar, baz: qux}` + "\n",
			location: "/baz",
			value:    [2]int{1, 17},
			property: [2]int{1, 12},
		},
		{
			test:     "FlowSequence",
			node:     `[1, 2, 3]` + "\n",
			location: "/1",
			value:    [2]int{1, 5},
			property: [2]int{1, 5},
		},
		{
			test: "ScalarRoot",
			node: heredoc.Doc(`
				hello
			`),
			location: "",
			value:    [2]int{1, 1},
			property: [2]int{1, 1},
		},
		{
			test: "ScalarOverTraverse",
			node: heredoc.Doc(`
				foo: bar
			`),
			location: "/foo/baz",
			value:    [2]int{0, 0},
			property: [2]int{0, 0},
		},
		{
			test: "SequenceIndexNotNumber",
			node: heredoc.Doc(`
				foo:
				  - 1
			`),
			location: "/foo/abc",
			value:    [2]int{0, 0},
			property: [2]int{0, 0},
		},
		{
			test: "SequenceItemMapping",
			node: heredoc.Doc(`
				- foo: 1
				- bar: 2
			`),
			location: "/1",
			value:    [2]int{2, 6},
			property: [2]int{2, 6},
		},
		{
			test: "BlockScalarLiteral",
			node: heredoc.Doc(`
				foo: |
				  hello
				  world
			`),
			location: "/foo",
			value:    [2]int{1, 6},
			property: [2]int{1, 1},
		},
		{
			test: "BlockScalarFolded",
			node: heredoc.Doc(`
				foo: >
				  hello
			`),
			location: "/foo",
			value:    [2]int{1, 6},
			property: [2]int{1, 1},
		},
		{
			test: "Anchor",
			node: heredoc.Doc(`
				foo: &a 1
				bar: 2
			`),
			location: "/foo",
			value:    [2]int{1, 6},
			property: [2]int{1, 1},
		},
		{
			test: "Alias",
			node: heredoc.Doc(`
				foo: &a 1
				bar: *a
			`),
			location: "/bar",
			value:    [2]int{2, 6},
			property: [2]int{2, 1},
		},
		{
			test: "QuotedKey",
			node: heredoc.Doc(`
				"foo bar": baz
			`),
			location: "/foo bar",
			value:    [2]int{1, 12},
			property: [2]int{1, 1},
		},
		{
			test: "QuotedValue",
			node: heredoc.Doc(`
				foo: "hello"
			`),
			location: "/foo",
			value:    [2]int{1, 6},
			property: [2]int{1, 1},
		},
		{
			test: "DocumentMarker",
			node: heredoc.Doc(`
				---
				foo: bar
			`),
			location: "/foo",
			value:    [2]int{2, 6},
			property: [2]int{2, 1},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			file, err := parser.ParseBytes([]byte(test.node), parser.ParseComments)
			s.Require().NoError(err)

			node := file.Docs[0].Body
			l := yamlvalidation.Locator{Node: node}

			// Value
			line, column := l.ValueAt(test.location)
			s.Equal(test.value[0], line, "value line not equal")
			s.Equal(test.value[1], column, "value column not equal")

			// Property
			line, column = l.PropertyAt(test.location)
			s.Equal(test.property[0], line, "property line not equal")
			s.Equal(test.property[1], column, "property column not equal")
		})
	}
}
