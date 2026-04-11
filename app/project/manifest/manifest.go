package manifest

import (
	_ "embed"
	"slices"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/yaml/parser"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

const filename = ".manala.yaml"

//go:embed template.yaml.tmpl
var _template string

func New() *Manifest {
	return &Manifest{
		config: &config{},
		vars:   map[string]any{},
	}
}

type Manifest struct {
	config *config
	vars   map[string]any
}

func (manifest *Manifest) Recipe() string {
	return manifest.config.Recipe
}

func (manifest *Manifest) Repository() string {
	return manifest.config.Repository
}

func (manifest *Manifest) Vars() map[string]any {
	return manifest.vars
}

func (manifest *Manifest) UnmarshalYAML(content []byte) error {
	// Parse content to node
	node, err := parser.Parse(content)
	if err != nil {
		return err
	}

	// Partition config & vars
	i := slices.IndexFunc(node.Values, func(node *ast.MappingValueNode) bool {
		return node.Key.String() == "manala"
	})
	if i == -1 {
		return parser.ErrorFrom(
			serrors.New("missing manala property"),
		)
	}

	configNode := node.Values[i].Value
	node.Values = slices.Concat(node.Values[:i], node.Values[i+1:])

	// Decode config
	if err = yaml.NodeToValue(configNode, manifest.config,
		yaml.Validator(configValidator{}),
		yaml.DisallowUnknownField(),
	); err != nil {
		return parser.ErrorFrom(err)
	}

	// Decode vars
	if err = yaml.NodeToValue(node, &manifest.vars); err != nil {
		return parser.ErrorFrom(err)
	}

	return nil
}
