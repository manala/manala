package yaml

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	goYamlParser "github.com/goccy/go-yaml/parser"
	"github.com/stretchr/testify/suite"
	"manala/internal/errors/serrors"
	"manala/internal/testing/heredoc"
	"testing"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestError() {
	_, formattedErr := goYamlParser.ParseBytes([]byte(`&foo`), 0)

	tests := []struct {
		test     string
		err      error
		expected *serrors.Assert
	}{
		{
			test: "Unknown",
			err:  serrors.New(`error`),
			expected: &serrors.Assert{
				Type:    &Error{},
				Message: "error",
			},
		},
		{
			test: "Formatted",
			err:  formattedErr,
			expected: &serrors.Assert{
				Type:    &Error{},
				Message: "unexpected anchor. anchor value is undefined",
				Arguments: []any{
					"line", 1,
					"column", 2,
				},
				Details: heredoc.Doc(`
					>  1 | &foo
					        ^
				`),
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

func (s *ErrorsSuite) TestNodeError() {
	contentNode, _ := NewParser().ParseBytes([]byte(`foo: bar`))

	tests := []struct {
		test     string
		node     goYamlAst.Node
		expected *serrors.Assert
	}{
		{
			test: "Content",
			node: contentNode,
			expected: &serrors.Assert{
				Type:    &NodeError{},
				Message: "error",
				Arguments: []any{
					"line", 1,
					"column", 4,
				},
				Details: heredoc.Doc(`
					>  1 | foo: bar
					          ^
				`),
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
