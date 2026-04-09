package yaml

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/printer"
)

type NodeTrace struct {
	Line        int
	Column      int
	DetailsFunc func(ansi bool) string
}

func NewNodeTrace(node ast.Node) NodeTrace {
	// Token
	token := node.GetToken()

	return NodeTrace{
		Line:   token.Position.Line,
		Column: token.Position.Column,
		DetailsFunc: func(ansi bool) string {
			var pp printer.Printer

			// Ensure there is *always* a trailing line feed
			return strings.TrimRight(
				pp.PrintErrorToken(token, ansi),
				"\n",
			) + "\n"
		},
	}
}
