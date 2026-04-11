package path_test

import (
	"testing"

	"github.com/manala/manala/internal/yaml/path"

	"github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type NodeSuite struct{ suite.Suite }

func TestNodeSuite(t *testing.T) {
	suite.Run(t, new(NodeSuite))
}

func (s *NodeSuite) TestNode() {
	node := ast.Null(nil)

	tests := []struct {
		test     string
		path     string
		expected string
	}{
		{
			test:     "Root",
			path:     "$",
			expected: "",
		},
		{
			test:     "FirstLevel",
			path:     "$.foo",
			expected: "foo",
		},
		{
			test:     "Levels",
			path:     "$.foo.bar",
			expected: "foo.bar",
		},
		{
			test:     "Index",
			path:     "$.foo[0].bar",
			expected: "foo[0].bar",
		},
		{
			test:     "IndexLast",
			path:     "$.foo[0]",
			expected: "foo[0]",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node.Path = test.path

			p := path.NewNodePath(node)

			s.Equal(test.expected, p.String())
		})
	}
}
