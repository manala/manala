package manifest

import (
	"slices"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/sync"
	"github.com/manala/manala/internal/yaml/parser"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

const filename = ".manala.yaml"

func New() *Manifest {
	return &Manifest{
		config: &config{
			Sync: recipe.Sync{},
		},
		vars:   map[string]any{},
		schema: schema.Schema{},
	}
}

type Manifest struct {
	config  *config
	vars    map[string]any
	schema  schema.Schema
	options []app.RecipeOption
}

func (manifest *Manifest) Description() string {
	return manifest.config.Description
}

func (manifest *Manifest) Icon() string {
	return manifest.config.Icon
}

func (manifest *Manifest) Template() string {
	return manifest.config.Template
}

func (manifest *Manifest) Partials() []string {
	return manifest.config.Partials
}

func (manifest *Manifest) Vars() map[string]any {
	return manifest.vars
}

func (manifest *Manifest) Sync() []sync.UnitInterface {
	return manifest.config.Sync
}

func (manifest *Manifest) Schema() schema.Schema {
	return manifest.schema
}

func (manifest *Manifest) Options() []app.RecipeOption {
	return manifest.options
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

	// Infer schema & options
	inf := Inferrer{
		Schema:  &manifest.schema,
		Options: &manifest.options,
	}
	if err = inf.Infer(node); err != nil {
		return err
	}

	return nil
}
