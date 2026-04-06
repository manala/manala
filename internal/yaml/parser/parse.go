package parser

import (
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/yaml"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

func Parse(bytes []byte) (*ast.MappingNode, error) {
	file, err := parser.ParseBytes(bytes, parser.ParseComments)
	if err != nil {
		return nil, yaml.NewError(err)
	}

	// File must not be empty...
	if len(file.Docs) == 0 || file.Docs[0].Body == nil {
		return nil, serrors.New("empty yaml file")
	}

	// ... nor include multiple documents
	if len(file.Docs) > 1 {
		return nil, yaml.NewNodeError("multiple documents yaml file", file.Docs[1].Body)
	}

	// ... and the first document must be a map
	node, ok := file.Docs[0].Body.(*ast.MappingNode)
	if !ok {
		return nil, yaml.NewNodeError("yaml document must be a map", file.Docs[0].Body)
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
