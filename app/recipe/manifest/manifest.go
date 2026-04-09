package manifest

import (
	"slices"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/sync"
	"github.com/manala/manala/internal/yaml"
	"github.com/manala/manala/internal/yaml/parser"

	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
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
	node    *goYamlAst.MappingNode
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
	var err error

	// Parse content to node
	manifest.node, err = parser.Parse(content)
	if err != nil {
		return err
	}

	// Partition config & vars
	i := slices.IndexFunc(manifest.node.Values, func(node *goYamlAst.MappingValueNode) bool {
		return node.Key.String() == "manala"
	})
	if i == -1 {
		return &parsing.Error{
			Err: serrors.New("missing manala property"),
		}
	}

	configNode := manifest.node.Values[i].Value
	manifest.node.Values = slices.Concat(manifest.node.Values[:i], manifest.node.Values[i+1:])

	// Decode config
	if err = goYaml.NodeToValue(configNode, manifest.config,
		goYaml.Validator(configValidator{}),
		goYaml.DisallowUnknownField(),
	); err != nil {
		return parser.ErrorFrom(err)
	}

	// Decode vars
	if err = goYaml.NodeToValue(manifest.node, &manifest.vars); err != nil {
		return parser.ErrorFrom(err)
	}

	// Infer schema
	if err = yaml.NewNodeSchemaInferrer(manifest.node).Infer(manifest.schema); err != nil {
		return serrors.New("unable to infer recipe manifest schema").
			WithErrors(err)
	}

	// Infer options
	if err = option.NewInferrer().Infer(manifest.node, &manifest.options); err != nil {
		return serrors.New("unable to infer recipe manifest options").
			WithErrors(err)
	}

	return nil
}
