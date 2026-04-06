package parser

import (
	"github.com/manala/manala/internal/yaml"

	"github.com/goccy/go-yaml/ast"
)

func resolve(node ast.Node, anchors map[string]ast.Node) error {
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
						return yaml.NewNodeError("cannot find anchor", vv.Value).
							WithArguments("anchor", alias)
					}

					switch a := anchor.(type) {
					case *ast.MappingNode:
						mergedValues = a.Values
					default:
						return yaml.NewNodeError("anchor must be a map", anchor).
							WithArguments("anchor", alias)
					}
				} else {
					return yaml.NewNodeError("map value must be an alias", v.Value)
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
				if err := resolveValue(&mv.Value, anchors); err != nil {
					return err
				}
			}
		}

		n.Values = deduplicatedValues

	case *ast.SequenceNode:
		for idx := range n.Values {
			if err := resolveValue(&n.Values[idx], anchors); err != nil {
				return err
			}
		}
	}

	return nil
}

func resolveValue(node *ast.Node, anchors map[string]ast.Node) error {
	switch n := (*node).(type) {
	case *ast.TagNode:
		*node = n.Value
		return resolveValue(node, anchors)
	case *ast.MappingKeyNode:
		*node = n.Value
		return resolveValue(node, anchors)
	case *ast.AliasNode:
		alias := n.Value.GetToken().Value
		anchor := anchors[alias]

		if anchor == nil {
			return yaml.NewNodeError("cannot find anchor", n.Value).
				WithArguments("anchor", alias)
		}

		*node = anchor
		return resolveValue(node, anchors)
	case *ast.AnchorNode:
		*node = n.Value
		return resolveValue(node, anchors)
	default:
		return resolve(*node, anchors)
	}
}
