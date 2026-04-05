package parser

import (
	"github.com/manala/manala/internal/yaml"

	"github.com/goccy/go-yaml/ast"
)

func resolve(node ast.Node, anchors map[string]ast.Node) (ast.Node, error) {
	switch n := node.(type) {
	case *ast.MappingNode:
		deduplicatedValues := make([]*ast.MappingValueNode, 0)

		for _, v := range n.Values {
			// Merge values
			mergedValues := make([]*ast.MappingValueNode, 0)

			if _, ok := v.Key.(*ast.MergeKeyNode); ok {
				if vv, ok := v.Value.(*ast.AliasNode); ok {
					alias := vv.Value.GetToken().Value

					anchor := anchors[alias]
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
				value, err := resolve(mv.Value, anchors)
				if err != nil {
					return nil, err
				}

				mv.Value = value
			}
		}

		n.Values = deduplicatedValues

		return n, nil
	case *ast.TagNode:
		return resolve(n.Value, anchors)
	case *ast.MappingKeyNode:
		return resolve(n.Value, anchors)
	case *ast.SequenceNode:
		for idx, v := range n.Values {
			value, err := resolve(v, anchors)
			if err != nil {
				return nil, err
			}

			n.Values[idx] = value
		}
	case *ast.AliasNode:
		alias := n.Value.GetToken().Value
		anchor := anchors[alias]

		if anchor == nil {
			return nil, yaml.NewNodeError("cannot find anchor", n.Value).
				WithArguments("anchor", alias)
		}

		return resolve(anchor, anchors)
	case *ast.AnchorNode:
		return resolve(n.Value, anchors)
	}

	return node, nil
}
