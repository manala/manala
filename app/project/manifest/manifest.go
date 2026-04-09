package manifest

import (
	_ "embed"
	"slices"

	"github.com/manala/manala/internal/parsing"
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
	node   *ast.MappingNode
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
	var err error

	// Parse content to node
	manifest.node, err = parser.Parse(content)
	if err != nil {
		return err
	}

	// Partition config & vars
	i := slices.IndexFunc(manifest.node.Values, func(node *ast.MappingValueNode) bool {
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
	if err = yaml.NodeToValue(configNode, manifest.config,
		yaml.Validator(configValidator{}),
		yaml.DisallowUnknownField(),
	); err != nil {
		return parser.ErrorFrom(err)
	}

	// Decode vars
	if err = yaml.NodeToValue(manifest.node, &manifest.vars); err != nil {
		return parser.ErrorFrom(err)
	}

	return nil
}
