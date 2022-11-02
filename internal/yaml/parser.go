package yaml

import (
	"fmt"
	yamlAst "github.com/goccy/go-yaml/ast"
	yamlParser "github.com/goccy/go-yaml/parser"
	"os"
)

func NewParser(opts ...ParserOption) *Parser {
	p := &Parser{
		anchors: map[string]yamlAst.Node{},
	}

	// Options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

type Parser struct {
	comments bool
	anchors  map[string]yamlAst.Node
	err      error
}

func (parser *Parser) ParseFile(filename string) (yamlAst.Node, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	node, err := parser.ParseBytes(content)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (parser *Parser) ParseBytes(bytes []byte) (yamlAst.Node, error) {
	// Parse with comments ?
	var mode yamlParser.Mode = 0
	if parser.comments {
		mode = yamlParser.ParseComments
	}

	file, err := yamlParser.ParseBytes(bytes, mode)
	if err != nil {
		return nil, NewError(err)
	}

	// File must not be empty...
	if len(file.Docs) == 0 {
		return nil, fmt.Errorf("empty yaml file")
	}

	// ... nor include multiple documents
	if len(file.Docs) > 1 {
		return nil, NewNodeError("multiple documents yaml file", file.Docs[1].Body)
	}

	node := file.Docs[0].Body

	yamlAst.Walk(parser, node)

	if parser.err != nil {
		return nil, parser.err
	}

	node, err = parser.resolve(node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (parser *Parser) Visit(node yamlAst.Node) yamlAst.Visitor {
	// Comment of the first MappingValueNode is set on its MappingNode.
	// Work around by manually move it.
	// See: https://github.com/goccy/go-yaml/issues/311
	if parser.comments {
		if n, ok := node.(*yamlAst.MappingNode); ok {
			if len(n.Values) > 0 && n.Comment != nil {
				n.Values[0].Comment = n.Comment
				n.Comment = nil
			}
		}
	}

	// Irregular map keys
	if n, ok := node.(*yamlAst.MappingValueNode); ok {
		if _, ok := n.Key.(*yamlAst.MergeKeyNode); ok {
			return parser
		}
		if _, ok := n.Key.(*yamlAst.StringNode); ok {
			return parser
		}
		parser.err = NewNodeError("irregular map key", node)
		return nil
	}

	switch n := node.(type) {
	case *yamlAst.AnchorNode:
		// Store anchors for coming resolution
		anchorName := n.Name.GetToken().Value
		parser.anchors[anchorName] = n.Value
	case
		// Scalars
		*yamlAst.NullNode,
		*yamlAst.IntegerNode,
		*yamlAst.FloatNode,
		*yamlAst.StringNode, *yamlAst.LiteralNode,
		*yamlAst.BoolNode,
		// Maps
		*yamlAst.MappingKeyNode,
		yamlAst.MapNode,
		// Arrays
		yamlAst.ArrayNode,
		// Aliases
		*yamlAst.AliasNode, *yamlAst.MergeKeyNode,
		// Tags
		*yamlAst.TagNode,
		// Comments
		*yamlAst.CommentGroupNode:
		// ¯\_(ツ)_/¯
	default:
		// Irregular types
		parser.err = NewNodeError("irregular type", node)
		return nil
	}

	return parser
}

func (parser *Parser) resolve(node yamlAst.Node) (yamlAst.Node, error) {
	switch n := node.(type) {
	case yamlAst.MapNode:
		values := make([]*yamlAst.MappingValueNode, 0)
		if m, ok := n.(*yamlAst.MappingNode); ok {
			values = m.Values
		} else {
			values = append(values, n.(*yamlAst.MappingValueNode))
		}

		deduplicatedValues := make([]*yamlAst.MappingValueNode, 0)

		for _, v := range values {
			// Merge values
			mergedValues := make([]*yamlAst.MappingValueNode, 0)
			if _, ok := v.Key.(*yamlAst.MergeKeyNode); ok {
				if vv, ok := v.Value.(*yamlAst.AliasNode); ok {
					alias := vv.Value.GetToken().Value
					anchor := parser.anchors[alias]
					if anchor == nil {
						return nil, NewNodeError(
							fmt.Sprintf("cannot find anchor \"%s\"", alias),
							vv.Value,
						)
					}
					switch a := anchor.(type) {
					case *yamlAst.MappingNode:
						mergedValues = a.Values
					case *yamlAst.MappingValueNode:
						mergedValues = append(mergedValues, a)
					default:
						return nil, NewNodeError(
							fmt.Sprintf("anchor \"%s\" must be a map", alias),
							anchor,
						)
					}
				} else {
					return nil, NewNodeError("map value must be an alias", v.Value)
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

		// Return either MappingValue or Mapping node,
		// depending on deduplicated values number
		if len(deduplicatedValues) == 1 {
			return deduplicatedValues[0], nil
		} else {
			if m, ok := n.(*yamlAst.MappingNode); ok {
				m.Values = deduplicatedValues
				return m, nil
			}

			m := &yamlAst.MappingNode{
				BaseNode: &yamlAst.BaseNode{},
			}
			m.Values = deduplicatedValues
			return m, nil
		}
	case *yamlAst.TagNode:
		return parser.resolve(n.Value)
	case *yamlAst.MappingKeyNode:
		return parser.resolve(n.Value)
	case *yamlAst.SequenceNode:
		for idx, v := range n.Values {
			value, err := parser.resolve(v)
			if err != nil {
				return nil, err
			}
			n.Values[idx] = value
		}
	case *yamlAst.AliasNode:
		alias := n.Value.GetToken().Value
		anchor := parser.anchors[alias]
		if anchor == nil {
			return nil, NewNodeError(
				fmt.Sprintf("cannot find anchor \"%s\"", alias),
				n.Value,
			)
		}
		return parser.resolve(anchor)
	case *yamlAst.AnchorNode:
		return parser.resolve(n.Value)
	}

	return node, nil
}

type ParserOption func(parser *Parser)

func WithComments() ParserOption {
	return func(parser *Parser) {
		parser.comments = true
	}
}
