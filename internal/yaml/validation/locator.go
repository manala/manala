package validation

import (
	"github.com/manala/manala/internal/validation"
	yamlpath "github.com/manala/manala/internal/yaml/path"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

func WithLocator(node ast.Node) validation.ValidateOption {
	return validation.WithLocator(Locator{Node: node})
}

// Locator resolves a JSON pointer (RFC 6901) to a line/column position within a YAML node.
type Locator struct {
	Node ast.Node
}

func (l Locator) ValueAt(location string) (int, int) {
	node := l.at(location)
	if node == nil {
		return 0, 0
	}

	if n, ok := node.(*ast.MappingNode); ok {
		parent := ast.Parent(l.Node, n)

		// Root map
		if parent == node {
			return 0, 0
		}

		if parent, ok := parent.(*ast.MappingValueNode); ok {
			node = parent.Key
		}
	}

	token := node.GetToken()
	if token == nil {
		return 0, 0
	}

	return token.Position.Line, token.Position.Column
}

func (l Locator) PropertyAt(location string) (int, int) {
	node := l.at(location)
	if node == nil {
		return 0, 0
	}

	parent := ast.Parent(l.Node, node)
	if parent, ok := parent.(*ast.MappingValueNode); ok {
		node = parent.Key
	}

	token := node.GetToken()
	if token == nil {
		return 0, 0
	}

	return token.Position.Line, token.Position.Column
}

func (l Locator) at(location string) ast.Node {
	if l.Node == nil {
		return nil
	}

	// Root
	if location == "" {
		return l.Node
	}

	path, err := yaml.PathString(yamlpath.FromJSONPointer(location))
	if err != nil {
		return nil
	}

	node, err := path.FilterNode(l.Node)
	if err != nil || node == nil {
		return nil
	}

	return node
}
