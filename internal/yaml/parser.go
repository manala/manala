package yaml

import (
	"os"

	"github.com/manala/manala/internal/serrors"

	goYamlAst "github.com/goccy/go-yaml/ast"
	goYamlParser "github.com/goccy/go-yaml/parser"
)

type Parser struct {
	comments bool
	anchors  map[string]goYamlAst.Node
	err      error
}

func NewParser(opts ...ParserOption) *Parser {
	p := &Parser{
		anchors: map[string]goYamlAst.Node{},
	}

	// Options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (parser *Parser) ParseFile(filename string) (goYamlAst.Node, error) {
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

func (parser *Parser) ParseBytes(bytes []byte) (goYamlAst.Node, error) {
	// Parse with comments ?
	var mode goYamlParser.Mode
	if parser.comments {
		mode = goYamlParser.ParseComments
	}

	file, err := goYamlParser.ParseBytes(bytes, mode)
	if err != nil {
		return nil, NewError(err)
	}

	// File must not be empty...
	if len(file.Docs) == 0 || file.Docs[0].Body == nil {
		return nil, serrors.New("empty yaml file")
	}

	// ... nor include multiple documents
	if len(file.Docs) > 1 {
		return nil, NewNodeError("multiple documents yaml file", file.Docs[1].Body)
	}

	node := file.Docs[0].Body

	goYamlAst.Walk(parser, node)

	if parser.err != nil {
		return nil, parser.err
	}

	node, err = parser.resolve(node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (parser *Parser) Visit(node goYamlAst.Node) goYamlAst.Visitor {
	switch n := node.(type) {
	case *goYamlAst.AnchorNode:
		// Store anchors for further resolution
		anchorName := n.Name.GetToken().Value
		parser.anchors[anchorName] = n.Value
		return parser
	case *goYamlAst.MappingValueNode:
		switch n.Key.(type) {
		case
			*goYamlAst.MergeKeyNode,
			*goYamlAst.StringNode,
			*goYamlAst.MappingKeyNode:
			return parser
		}

		parser.err = NewNodeError("irregular map key", n)
		return nil
	case *goYamlAst.MappingKeyNode:
		switch n.Value.(type) {
		case
			*goYamlAst.MergeKeyNode,
			*goYamlAst.StringNode:
			return parser
		}

		parser.err = NewNodeError("irregular map key", n)
		return nil
	case
		// Scalars
		*goYamlAst.NullNode,
		*goYamlAst.IntegerNode,
		*goYamlAst.FloatNode,
		*goYamlAst.StringNode, *goYamlAst.LiteralNode,
		*goYamlAst.BoolNode,
		// Maps
		goYamlAst.MapNode,
		// Tags
		*goYamlAst.TagNode,
		// Arrays
		goYamlAst.ArrayNode,
		// Aliases
		*goYamlAst.AliasNode, *goYamlAst.MergeKeyNode,
		// Comments
		*goYamlAst.CommentGroupNode:
		return parser
	}

	parser.err = NewNodeError("irregular type", node)
	return nil
}

func (parser *Parser) resolve(node goYamlAst.Node) (goYamlAst.Node, error) {
	switch n := node.(type) {
	case *goYamlAst.MappingNode:
		deduplicatedValues := make([]*goYamlAst.MappingValueNode, 0)

		for _, v := range n.Values {
			// Merge values
			mergedValues := make([]*goYamlAst.MappingValueNode, 0)

			if _, ok := v.Key.(*goYamlAst.MergeKeyNode); ok {
				if vv, ok := v.Value.(*goYamlAst.AliasNode); ok {
					alias := vv.Value.GetToken().Value

					anchor := parser.anchors[alias]
					if anchor == nil {
						return nil, NewNodeError("cannot find anchor", vv.Value).
							WithArguments("anchor", alias)
					}

					switch a := anchor.(type) {
					case *goYamlAst.MappingNode:
						mergedValues = a.Values
					default:
						return nil, NewNodeError("anchor must be a map", anchor).
							WithArguments("anchor", alias)
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

		n.Values = deduplicatedValues

		return n, nil
	case *goYamlAst.TagNode:
		return parser.resolve(n.Value)
	case *goYamlAst.MappingKeyNode:
		return parser.resolve(n.Value)
	case *goYamlAst.SequenceNode:
		for idx, v := range n.Values {
			value, err := parser.resolve(v)
			if err != nil {
				return nil, err
			}

			n.Values[idx] = value
		}
	case *goYamlAst.AliasNode:
		alias := n.Value.GetToken().Value
		anchor := parser.anchors[alias]

		if anchor == nil {
			return nil, NewNodeError("cannot find anchor", n.Value).
				WithArguments("anchor", alias)
		}

		return parser.resolve(anchor)
	case *goYamlAst.AnchorNode:
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
