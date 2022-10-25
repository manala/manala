package yaml

import (
	yamlAst "github.com/goccy/go-yaml/ast"
	yamlPrinter "github.com/goccy/go-yaml/printer"
	"github.com/muesli/termenv"
	internalReport "manala/internal/report"
)

func NewReporter(node yamlAst.Node) *Reporter {
	return &Reporter{
		node: node,
	}
}

type Reporter struct {
	node yamlAst.Node
}

func (reporter *Reporter) Report(report *internalReport.Report) {
	var pp yamlPrinter.Printer

	color := true
	if termenv.EnvColorProfile() == termenv.Ascii {
		color = false
	}

	if reporter.node != nil {
		token := reporter.node.GetToken()
		report.Compose(
			internalReport.WithField("line", token.Position.Line),
			internalReport.WithField("column", token.Position.Column),
			internalReport.WithTrace(pp.PrintErrorToken(token, color)),
		)
	}
}
