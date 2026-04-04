package parser

import (
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/yaml"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

type Parser struct {
	anchors map[string]ast.Node
	err     error
}

func NewParser() *Parser {
	return &Parser{
		anchors: map[string]ast.Node{},
	}
}

func (p *Parser) ParseBytes(bytes []byte) (ast.Node, error) {
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

	node := file.Docs[0].Body

	ast.Walk(p, node)

	if p.err != nil {
		return nil, p.err
	}

	node, err = p.resolve(node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (parser *Parser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.AnchorNode:
		// Store anchors for further resolution
		anchorName := n.Name.GetToken().Value
		parser.anchors[anchorName] = n.Value
		return parser
	case *ast.MappingValueNode:
		switch n.Key.(type) {
		case
			*ast.MergeKeyNode,
			*ast.StringNode,
			*ast.MappingKeyNode:
			return parser
		}

		parser.err = yaml.NewNodeError("irregular map key", n)
		return nil
	case *ast.MappingKeyNode:
		switch n.Value.(type) {
		case
			*ast.MergeKeyNode,
			*ast.StringNode:
			return parser
		}

		parser.err = yaml.NewNodeError("irregular map key", n)
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
		return parser
	}

	parser.err = yaml.NewNodeError("irregular type", node)
	return nil
}

func (parser *Parser) resolve(node ast.Node) (ast.Node, error) {
	switch n := node.(type) {
	case *ast.MappingNode:
		deduplicatedValues := make([]*ast.MappingValueNode, 0)

		for _, v := range n.Values {
			// Merge values
			mergedValues := make([]*ast.MappingValueNode, 0)

			if _, ok := v.Key.(*ast.MergeKeyNode); ok {
				if vv, ok := v.Value.(*ast.AliasNode); ok {
					alias := vv.Value.GetToken().Value

					anchor := parser.anchors[alias]
					if anchor == nil {
						return nil, yaml.NewNodeError("cannot find anchor", vv.Value).
							WithArguments("anchor", alias)
					}

					switch a := anchor.(type) {
					case *ast.MappingNode:
						mergedValues = a.Values
					default:
						return nil, yaml.NewNodeError("anchor must be a map", anchor).
							WithArguments("anchor", alias)
					}
				} else {
					return nil, yaml.NewNodeError("map value must be an alias", v.Value)
				}
			} else {
				mergedValues = append(mergedValues, v)
			}

			// Deduplicate values
			for _, mv := range mergedValues {
				for i, dv := range deduplicatedValues {
					if mv.Key.GetToken().Value == dv.Key.GetToken().Value {
						deduplicatedValues = append(deduplicatedValues[:i], deduplicatedValues[i+1:]...)

						break
					}
				}

				deduplicatedValues = append(deduplicatedValues, mv)

				// Resolve
				value, err := parser.resolve(mv.Value)
				if err != nil {
					return nil, err
				}

				mv.Value = value
			}
		}

		n.Values = deduplicatedValues

		return n, nil
	case *ast.TagNode:
		return parser.resolve(n.Value)
	case *ast.MappingKeyNode:
		return parser.resolve(n.Value)
	case *ast.SequenceNode:
		for idx, v := range n.Values {
			value, err := parser.resolve(v)
			if err != nil {
				return nil, err
			}

			n.Values[idx] = value
		}
	case *ast.AliasNode:
		alias := n.Value.GetToken().Value
		anchor := parser.anchors[alias]

		if anchor == nil {
			return nil, yaml.NewNodeError("cannot find anchor", n.Value).
				WithArguments("anchor", alias)
		}

		return parser.resolve(anchor)
	case *ast.AnchorNode:
		return parser.resolve(n.Value)
	}

	return node, nil
}
