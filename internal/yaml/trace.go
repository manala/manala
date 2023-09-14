package yaml

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	goYamlPrinter "github.com/goccy/go-yaml/printer"
	"strings"
)

func NewNodeTrace(node goYamlAst.Node) NodeTrace {
	// Token
	token := node.GetToken()

	return NodeTrace{
		Line:   token.Position.Line,
		Column: token.Position.Column,
		DetailsFunc: func(ansi bool) string {
			var pp goYamlPrinter.Printer

			// Ensure there is *always* a trailing line feed
			return strings.TrimRight(
				pp.PrintErrorToken(token, ansi),
				"\n",
			) + "\n"
		},
	}
}

type NodeTrace struct {
	Line        int
	Column      int
	DetailsFunc func(ansi bool) string
}
