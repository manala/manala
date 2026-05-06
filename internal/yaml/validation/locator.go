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

func (l Locator) At(location string) (int, int) {
	path, err := yaml.PathString(yamlpath.FromJSONPointer(location))
	if err != nil {
		return 0, 0
	}

	target, err := path.FilterNode(l.Node)
	if err != nil || target == nil {
		return 0, 0
	}

	token := target.GetToken()
	if token == nil {
		return 0, 0
	}

	return token.Position.Line, token.Position.Column
}
