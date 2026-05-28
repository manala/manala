package parser

import (
	"errors"
	"fmt"

	yamlerrors "github.com/manala/manala/internal/yaml/errors"

	"github.com/goccy/go-yaml/ast"
)

// resolve replaces aliases with their anchor values and deduplicates mapping keys.
// visiting tracks anchor names currently being resolved to detect cycles.
func resolve(node ast.Node, anchors map[string]ast.Node, visiting map[string]bool) error {
	switch n := node.(type) {
	case *ast.MappingNode:
		deduplicatedValues := make([]*ast.MappingValueNode, 0)

		for _, v := range n.Values {
			// Merge values
			mergedValues := make([]*ast.MappingValueNode, 0)

			if _, ok := v.Key.(*ast.MergeKeyNode); ok {
				switch vv := v.Value.(type) {
				case *ast.AliasNode:
					mn, err := resolveMergeAlias(vv, anchors, visiting)
					if err != nil {
						return err
					}
					mergedValues = mn.Values
				case *ast.SequenceNode:
					// `<<: [*a, *b]` — defer flattening to goccy's SequenceMergeValue
					// so we follow its iteration order rather than reimplementing it.
					maps := make([]ast.MapNode, 0, len(vv.Values))
					for _, elt := range vv.Values {
						alias, ok := elt.(*ast.AliasNode)
						if !ok {
							return yamlerrors.New(
								errors.New("map value must be an alias"),
								elt.GetToken(),
							)
						}
						mn, err := resolveMergeAlias(alias, anchors, visiting)
						if err != nil {
							return err
						}
						maps = append(maps, mn)
					}
					iter := ast.SequenceMergeValue(maps...).MapRange()
					for iter.Next() {
						mergedValues = append(mergedValues, iter.KeyValue())
					}
				default:
					return yamlerrors.New(
						errors.New("map value must be an alias"),
						v.Value.GetToken(),
					)
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
				if err := resolveValue(&mv.Value, anchors, visiting); err != nil {
					return err
				}
			}
		}

		n.Values = deduplicatedValues

	case *ast.SequenceNode:
		for idx := range n.Values {
			if err := resolveValue(&n.Values[idx], anchors, visiting); err != nil {
				return err
			}
		}
	}

	return nil
}

func resolveMergeAlias(alias *ast.AliasNode, anchors map[string]ast.Node, visiting map[string]bool) (*ast.MappingNode, error) {
	name := alias.Value.GetToken().Value

	if visiting[name] {
		return nil, yamlerrors.New(
			fmt.Errorf("cycle through yaml anchor %q", name),
			alias.GetToken(),
		)
	}

	anchor := anchors[name]
	if anchor == nil {
		return nil, yamlerrors.New(
			fmt.Errorf("unknown \"%s\" yaml anchor", name),
			alias.GetToken(),
		)
	}

	mn, ok := anchor.(*ast.MappingNode)
	if !ok {
		return nil, yamlerrors.New(
			fmt.Errorf("anchor %s must be a map", name),
			anchor.GetToken(),
		)
	}

	return mn, nil
}

func resolveValue(node *ast.Node, anchors map[string]ast.Node, visiting map[string]bool) error {
	switch n := (*node).(type) {
	case *ast.TagNode:
		*node = n.Value
		return resolveValue(node, anchors, visiting)
	case *ast.MappingKeyNode:
		*node = n.Value
		return resolveValue(node, anchors, visiting)
	case *ast.AliasNode:
		name := n.Value.GetToken().Value

		if visiting[name] {
			return yamlerrors.New(
				fmt.Errorf("cycle through yaml anchor %q", name),
				n.GetToken(),
			)
		}

		anchor := anchors[name]
		if anchor == nil {
			return yamlerrors.New(
				fmt.Errorf("unknown \"%s\" yaml anchor", name),
				n.GetToken(),
			)
		}

		visiting[name] = true
		*node = anchor
		err := resolveValue(node, anchors, visiting)
		delete(visiting, name)
		return err
	case *ast.AnchorNode:
		name := n.Name.GetToken().Value

		if visiting[name] {
			return yamlerrors.New(
				fmt.Errorf("cycle through yaml anchor %q", name),
				n.GetToken(),
			)
		}

		visiting[name] = true
		*node = n.Value
		err := resolveValue(node, anchors, visiting)
		delete(visiting, name)
		return err
	default:
		return resolve(*node, anchors, visiting)
	}
}
