package yaml_test

import (
	"testing"

	"github.com/manala/manala/internal/yaml"

	"github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type PathSuite struct{ suite.Suite }

func TestPathSuite(t *testing.T) {
	suite.Run(t, new(PathSuite))
}

func (s *PathSuite) TestNode() {
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

			path := yaml.NewNodePath(node)

			s.Equal(test.expected, path.String())
		})
	}
}
