package parser

import (
	"github.com/manala/manala/internal/yaml"

	"github.com/goccy/go-yaml/ast"
)

type walker struct {
	anchors map[string]ast.Node
	err     error
}

func (w *walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.AnchorNode:
		// Store anchors for further resolution
		anchorName := n.Name.GetToken().Value
		w.anchors[anchorName] = n.Value
		return w
	case *ast.MappingValueNode:
		switch n.Key.(type) {
		case
			*ast.MergeKeyNode,
			*ast.StringNode,
			*ast.MappingKeyNode:
			return w
		}

		w.err = yaml.NewNodeError("irregular map key", n)
		return nil
	case *ast.MappingKeyNode:
		switch n.Value.(type) {
		case
			*ast.MergeKeyNode,
			*ast.StringNode:
			return w
		}

		w.err = yaml.NewNodeError("irregular map key", n)
		return nil
	case
		// Scalars
		*ast.NullNode,
		*ast.IntegerNode,
		*ast.FloatNode,
		*ast.StringNode, *ast.LiteralNode,
		*ast.BoolNode,
		// Maps
		ast.MapNode,
		// Tags
		*ast.TagNode,
		// Arrays
		ast.ArrayNode,
		// Aliases
		*ast.AliasNode, *ast.MergeKeyNode,
		// Comments
		*ast.CommentGroupNode:
		return w
	}

	w.err = yaml.NewNodeError("irregular type", node)
	return nil
}
