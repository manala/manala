package parser

import (
	"errors"

	yamlerrors "github.com/manala/manala/internal/yaml/errors"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// Parse parses YAML bytes into a validated and resolved MappingNode,
// and returns an enhanced error with position information if parsing fails.
func Parse(data []byte) (*ast.MappingNode, error) {
	file, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return nil, yamlerrors.From(err)
	}

	// File must not be empty...
	if len(file.Docs) == 0 || file.Docs[0].Body == nil {
		return nil, yamlerrors.New(
			errors.New("empty yaml content"),
			nil,
		)
	}

	// ... nor include multiple documents
	if len(file.Docs) > 1 {
		return nil, yamlerrors.New(
			errors.New("multiple documents yaml content"),
			file.Docs[1].Start,
		)
	}

	// ... and the first document must be a map
	node, ok := file.Docs[0].Body.(*ast.MappingNode)
	if !ok {
		return nil, yamlerrors.New(
			errors.New("yaml document must be a map"),
			file.Docs[0].Body.GetToken(),
		)
	}

	// Walk
	w := &walker{
		anchors: map[string]ast.Node{},
	}
	ast.Walk(w, node)
	if w.err != nil {
		return nil, w.err
	}

	// Resolve
	if err := resolve(node, w.anchors, map[string]bool{}); err != nil {
		return nil, err
	}

	return node, nil
}
