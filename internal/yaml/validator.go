package yaml

import (
	"manala/internal/validator"

	goYamlAst "github.com/goccy/go-yaml/ast"
)

// NodeValidatorFormatter creates a node validator formatter.
func NodeValidatorFormatter(node goYamlAst.Node) validator.Formatter {
	return nodeValidatorFormatter{
		node: node,
	}
}

type nodeValidatorFormatter struct {
	node goYamlAst.Node
}

func (formatter nodeValidatorFormatter) Format(violation *validator.Violation) {
	// Get node by path
	node, err := NewNodePathAccessor(violation.Path).
		Get(formatter.node)
	if err != nil {
		return
	}

	// Trace
	trace := NewNodeTrace(node)
	violation.Line = trace.Line
	violation.Column = trace.Column
	violation.DetailsFunc = trace.DetailsFunc
}
