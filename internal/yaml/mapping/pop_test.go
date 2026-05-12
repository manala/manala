package mapping_test

import (
	"testing"

	"github.com/manala/manala/internal/testing/heredoc"
	yamlmapping "github.com/manala/manala/internal/yaml/mapping"
	yamlparser "github.com/manala/manala/internal/yaml/parser"

	"github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type PopSuite struct{ suite.Suite }

func TestPopSuite(t *testing.T) {
	suite.Run(t, new(PopSuite))
}

func (s *PopSuite) TestNotFound() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		foo: bar
	`)))
	s.Require().NoError(err)

	value, ok := yamlmapping.Pop(node, "baz")

	s.False(ok)
	s.Nil(value)
	s.Len(node.Values, 1)
}

func (s *PopSuite) TestFound() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		foo: bar
		baz: qux
	`)))
	s.Require().NoError(err)

	value, ok := yamlmapping.Pop(node, "foo")

	s.True(ok)
	s.Require().IsType((*ast.StringNode)(nil), value)
	s.Equal("bar", value.(*ast.StringNode).Value)
	s.Len(node.Values, 1)
	s.Equal("baz", node.Values[0].Key.String())
}

func (s *PopSuite) TestFoundLast() {
	node, err := yamlparser.Parse([]byte(heredoc.Doc(`
		foo: bar
		baz: qux
	`)))
	s.Require().NoError(err)

	value, ok := yamlmapping.Pop(node, "baz")

	s.True(ok)
	s.Require().IsType((*ast.StringNode)(nil), value)
	s.Equal("qux", value.(*ast.StringNode).Value)
	s.Len(node.Values, 1)
	s.Equal("foo", node.Values[0].Key.String())
}
