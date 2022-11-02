package recipe

import (
	_ "embed"
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"manala/core"
	internalReport "manala/internal/report"
	internalSyncer "manala/internal/syncer"
	internalValidation "manala/internal/validation"
	internalYaml "manala/internal/yaml"
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
	node    yamlAst.Node
	config  *manifestConfig
	vars    map[string]interface{}
	schema  map[string]interface{}
	options []core.RecipeOption
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

func (man *Manifest) Sync() []internalSyncer.UnitInterface {
	return man.config.Sync
}

func (man *Manifest) Schema() map[string]interface{} {
	return man.schema
}

func (man *Manifest) ReadFrom(reader io.Reader) error {
	// Read content
	content, err := io.ReadAll(reader)
	if err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to read recipe manifest")
	}

	// Parse content to node
	man.node, err = internalYaml.NewParser(internalYaml.WithComments()).ParseBytes(content)
	if err != nil {
		return internalReport.NewError(err).
			WithMessage("irregular recipe manifest")
	}

	// Validate node
	validation, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(manifestSchema),
		internalYaml.NewJsonLoader(man.node),
	)
	if err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to validate recipe manifest")
	}

	if !validation.Valid() {
		return internalReport.NewError(
			internalValidation.NewError(
				"invalid recipe manifest",
				validation,
			).
				WithReporter(man).
				WithMessages([]internalValidation.ErrorMessage{
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
				}),
		)
	}

	// Extract config node
	configNode, err := internalYaml.NewExtractor(&man.node).ExtractRootMap("manala")
	if err != nil {
		return internalReport.NewError(err).
			WithMessage("incorrect recipe manifest")
	}

	// Decode config
	if err = yaml.NodeToValue(configNode, man.config); err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to decode recipe manifest config")
	}

	// Decode vars
	if err = yaml.NodeToValue(man.node, &man.vars); err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to decode recipe manifest vars")
	}

	// Infer schema
	if err := NewSchemaInferrer().Infer(man.node, man.schema); err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to infer recipe manifest schema")
	}

	// Infer options
	if err := NewOptionsInferrer().Infer(man.node, &man.options); err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to infer recipe manifest options")
	}

	return nil
}

func (man *Manifest) Report(result gojsonschema.ResultError, report *internalReport.Report) {
	internalYaml.NewValidationPathReporter(man.node).Report(result, report)
}

func (man *Manifest) InitVars(callback func(options []core.RecipeOption) error) (map[string]interface{}, error) {
	var vars map[string]interface{}

	if err := callback(man.options); err != nil {
		return nil, internalReport.NewError(err).
			WithMessage("unable to apply recipe manifest options")
	}

	// Decode vars
	if err := yaml.NewDecoder(man.node).Decode(&vars); err != nil {
		return nil, internalReport.NewError(err).
			WithMessage("unable to decode recipe manifest init vars")
	}

	return vars, nil
}

type manifestConfig struct {
	Description string
	Template    string
	Sync        sync
}
