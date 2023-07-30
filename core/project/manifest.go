package project

import (
	_ "embed"
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"manala/internal/errors/serrors"
	"manala/internal/validation"
	"manala/internal/yaml"
)

//go:embed resources/manifest.schema.json
var manifestSchema string

func NewManifest() *Manifest {
	return &Manifest{
		config: &manifestConfig{},
		vars:   map[string]interface{}{},
	}
}

type Manifest struct {
	node   goYamlAst.Node
	config *manifestConfig
	vars   map[string]interface{}
}

func (man *Manifest) Recipe() string {
	return man.config.Recipe
}

func (man *Manifest) Repository() string {
	return man.config.Repository
}

func (man *Manifest) Vars() map[string]interface{} {
	return man.vars
}

func (man *Manifest) ReadFrom(reader io.Reader) error {
	// Read content
	content, err := io.ReadAll(reader)
	if err != nil {
		return serrors.Wrap("unable to read project manifest", err)
	}

	// Parse content to node
	man.node, err = yaml.NewParser().ParseBytes(content)
	if err != nil {
		return serrors.Wrap("irregular project manifest", err)
	}

	// Validate node
	val, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(manifestSchema),
		yaml.NewJsonLoader(man.node),
	)
	if err != nil {
		return serrors.Wrap("unable to validate project manifest", err)
	}

	if !val.Valid() {
		return validation.NewError(
			"invalid project manifest",
			val,
		).
			WithResultErrorDecorator(man.ValidationResultErrorDecorator()).
			WithMessages([]validation.ErrorMessage{
				{Field: "(root)", Type: "invalid_type", Message: "yaml document must be a map"},
				{Field: "(root)", Type: "required", Property: "manala", Message: "missing manala field"},
				{Field: "manala", Type: "invalid_type", Message: "manala field must be a map"},
				{Field: "manala", Type: "required", Property: "recipe", Message: "missing manala recipe field"},
				{Field: "manala", Type: "additional_property_not_allowed", Message: "manala field don't support additional properties"},
				// Recipe
				{Field: "manala.recipe", Type: "invalid_type", Message: "manala recipe field must be a string"},
				{Field: "manala.recipe", Type: "string_gte", Message: "empty manala recipe field"},
				{Field: "manala.recipe", Type: "string_lte", Message: "too long manala recipe field"},
				// Repository
				{Field: "manala.repository", Type: "invalid_type", Message: "manala repository field must be a string"},
				{Field: "manala.repository", Type: "string_gte", Message: "empty manala repository field"},
				{Field: "manala.repository", Type: "string_lte", Message: "too long manala repository field"},
			})
	}

	// Extract config node
	configNode, err := yaml.NewExtractor(&man.node).ExtractRootMap("manala")
	if err != nil {
		return serrors.Wrap("incorrect project manifest", err)
	}

	// Decode config
	if err = goYaml.NodeToValue(configNode, man.config); err != nil {
		return serrors.Wrap("unable to decode project manifest config", err)
	}

	// Decode vars
	if err = goYaml.NodeToValue(man.node, &man.vars); err != nil {
		return serrors.Wrap("unable to decode recipe manifest vars", err)
	}

	return nil
}

func (man *Manifest) ValidationResultErrorDecorator() validation.ResultErrorDecorator {
	return yaml.NewNodeValidationResultPathErrorDecorator(man.node)
}

type manifestConfig struct {
	Recipe     string
	Repository string
}
