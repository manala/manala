package yaml

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	goYamlParser "github.com/goccy/go-yaml/parser"
	"manala/internal/serrors"
)

func (s *Suite) TestError() {
	_, err := goYamlParser.ParseBytes([]byte(`&foo`), 0)

	tests := []struct {
		test     string
		err      error
		expected *serrors.Assert
	}{
		{
			test: "Unknown",
			err:  serrors.New("error"),
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "error",
			},
		},
		{
			test: "Formatted",
			err:  err,
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "unexpected anchor. anchor value is undefined",
				Arguments: []any{
					"line", 1,
					"column", 2,
				},
				Details: `
					>  1 | &foo
					        ^
				`,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := NewError(test.err)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestNodeError() {
	node, _ := NewParser().ParseBytes([]byte(`foo: bar`))

	tests := []struct {
		test     string
		node     goYamlAst.Node
		expected *serrors.Assert
	}{
		{
			test: "Content",
			node: node,
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			err := NewNodeError("error", test.node)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}
