package mapping

import (
	"slices"

	"github.com/goccy/go-yaml/ast"
)

// Pop removes the entry with the given key from the mapping node and returns its value.
// Returns (nil, false) if the key is not found.
func Pop(mapping *ast.MappingNode, key string) (ast.Node, bool) {
	i := slices.IndexFunc(mapping.Values, func(v *ast.MappingValueNode) bool {
		return v.Key.String() == key
	})
	if i == -1 {
		return nil, false
	}

	node := mapping.Values[i].Value
	mapping.Values = slices.Concat(mapping.Values[:i], mapping.Values[i+1:])

	return node, true
}
