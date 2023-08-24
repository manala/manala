package recipe

import (
	_ "embed"
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"manala/app/interfaces"
	"manala/internal/errors/serrors"
	"manala/internal/syncer"
	"manala/internal/validation"
	"manala/internal/yaml"
	"regexp"
)

//go:embed resources/manifest.schema.json
var manifestSchema string

func NewManifest() *Manifest {
	return &Manifest{
		config: &manifestConfig{
			Sync: sync{},
		},
		vars:   map[string]interface{}{},
		schema: map[string]interface{}{},
	}
}

type Manifest struct {
	node    goYamlAst.Node
	config  *manifestConfig
	vars    map[string]interface{}
	schema  map[string]interface{}
	options []interfaces.RecipeOption
}

func (man *Manifest) Description() string {
	return man.config.Description
}

func (man *Manifest) Template() string {
	return man.config.Template
}

func (man *Manifest) Vars() map[string]interface{} {
	return man.vars
}

func (man *Manifest) Sync() []syncer.UnitInterface {
	return man.config.Sync
}

func (man *Manifest) Schema() map[string]interface{} {
	return man.schema
}

func (man *Manifest) ReadFrom(reader io.Reader) error {
	// Read content
	content, err := io.ReadAll(reader)
	if err != nil {
		return serrors.Wrap("unable to read recipe manifest", err)
	}

	// Parse content to node
	man.node, err = yaml.NewParser(yaml.WithComments()).ParseBytes(content)
	if err != nil {
		return serrors.Wrap("irregular recipe manifest", err)
	}

	// Validate node
	val, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(manifestSchema),
		yaml.NewJsonLoader(man.node),
	)
	if err != nil {
		return serrors.Wrap("unable to validate recipe manifest", err)
	}

	if !val.Valid() {
		return validation.NewError(
			"invalid recipe manifest",
			val,
		).
			WithResultErrorDecorator(man.ValidationResultErrorDecorator()).
			WithMessages([]validation.ErrorMessage{
				{Field: "(root)", Type: "invalid_type", Message: "yaml document must be a map"},
				{Field: "(root)", Type: "required", Property: "manala", Message: "missing manala field"},
				{Field: "manala", Type: "invalid_type", Message: "manala field must be a map"},
				{Field: "manala", Type: "required", Property: "description", Message: "missing manala description field"},
				{Field: "manala", Type: "additional_property_not_allowed", Message: "manala field don't support additional properties"},
				// Description
				{Field: "manala.description", Type: "invalid_type", Message: "manala description field must be a string"},
				{Field: "manala.description", Type: "string_gte", Message: "empty manala description field"},
				{Field: "manala.description", Type: "string_lte", Message: "too long manala description field"},
				// Template
				{Field: "manala.template", Type: "invalid_type", Message: "manala template field must be a string"},
				{Field: "manala.template", Type: "string_gte", Message: "empty manala template field"},
				{Field: "manala.template", Type: "string_lte", Message: "too long manala template field"},
				// Sync
				{Field: "manala.sync", Type: "invalid_type", Message: "manala sync field must be a sequence"},
				// Sync Item
				{FieldRegex: regexp.MustCompile(`manala\.sync\.\d+`), Type: "invalid_type", Message: "manala sync sequence entries must be strings"},
				{FieldRegex: regexp.MustCompile(`manala\.sync\.\d+`), Type: "string_gte", Message: "empty manala sync sequence entry"},
				{FieldRegex: regexp.MustCompile(`manala\.sync\.\d+`), Type: "string_lte", Message: "too long manala sync sequence entry"},
			})
	}

	// Extract config node
	configNode, err := yaml.NewExtractor(&man.node).ExtractRootMap("manala")
	if err != nil {
		return serrors.Wrap("incorrect recipe manifest", err)
	}

	// Decode config
	if err = goYaml.NodeToValue(configNode, man.config); err != nil {
		return serrors.Wrap("unable to decode recipe manifest config", err)
	}

	// Decode vars
	if err = goYaml.NodeToValue(man.node, &man.vars); err != nil {
		return serrors.Wrap("unable to decode recipe manifest vars", err)
	}

	// Infer schema
	if err := NewSchemaInferrer().Infer(man.node, man.schema); err != nil {
		return serrors.Wrap("unable to infer recipe manifest schema", err)
	}

	// Infer options
	if err := NewOptionsInferrer().Infer(man.node, &man.options); err != nil {
		return serrors.Wrap("unable to infer recipe manifest options", err)
	}

	return nil
}

func (man *Manifest) ValidationResultErrorDecorator() validation.ResultErrorDecorator {
	return yaml.NewNodeValidationResultPathErrorDecorator(man.node)
}

func (man *Manifest) InitVars(callback func(options []interfaces.RecipeOption) error) (map[string]interface{}, error) {
	var vars map[string]interface{}

	if err := callback(man.options); err != nil {
		return nil, serrors.Wrap("unable to apply recipe manifest options", err)
	}

	// Decode vars
	if err := goYaml.NewDecoder(man.node).Decode(&vars); err != nil {
		return nil, serrors.Wrap("unable to decode recipe manifest init vars", err)
	}

	return vars, nil
}

type manifestConfig struct {
	Description string
	Template    string
	Sync        sync
}
