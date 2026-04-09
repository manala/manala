package parser

import (
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// Parse parses yaml bytes into a validated and resolved MappingNode.
func Parse(data []byte) (*ast.MappingNode, error) {
	file, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return nil, ErrorFrom(err)
	}

	// File must not be empty...
	if len(file.Docs) == 0 || file.Docs[0].Body == nil {
		return nil, &parsing.Error{
			Err: serrors.New("empty yaml content"),
		}
	}

	// ... nor include multiple documents
	if len(file.Docs) > 1 {
		return nil, ErrorAt(
			serrors.New("multiple documents yaml content"),
			file.Docs[1].GetToken().Prev.Prev,
		)
	}

	// ... and the first document must be a map
	node, ok := file.Docs[0].Body.(*ast.MappingNode)
	if !ok {
		return nil, ErrorAt(
			serrors.New("yaml document must be a map"),
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
	if err := resolve(node, w.anchors); err != nil {
		return nil, err
	}

	return node, nil
}
