package yaml

import (
	"github.com/manala/manala/internal/serrors"

	"github.com/goccy/go-yaml/ast"
)

func NewNodeError(message string, node ast.Node) serrors.Error {
	err := serrors.New(message)

	if node == nil {
		return err
	}

	// Trace
	trace := NewNodeTrace(node)

	return err.
		WithArguments(
			"line", trace.Line,
			"column", trace.Column,
		).
		WithDetailsFunc(trace.DetailsFunc)
}
