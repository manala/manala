package yaml_test

import (
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/yaml"
	"github.com/manala/manala/internal/yaml/parser"

	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestNode() {
	node, _ := parser.Parse([]byte(`foo: bar`))

	tests := []struct {
		test     string
		node     goYamlAst.Node
		expected errors.Assertion
	}{
		{
			test: "Content",
			node: node,
			expected: &serrors.Assertion{
				Message: "error",
				Arguments: []any{
					"line", 1,
					"column", 4,
				},
				Details: `
					>  1 | foo: bar
					          ^
				`,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := yaml.NewNodeError("error", test.node)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}
