package project

import (
	_ "embed"
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/xeipuuv/gojsonschema"
	"io"
	internalReport "manala/internal/report"
	internalValidation "manala/internal/validation"
	internalYaml "manala/internal/yaml"
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
	node   yamlAst.Node
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
		return internalReport.NewError(err).
			WithMessage("unable to read project manifest")
	}

	// Parse content to node
	man.node, err = internalYaml.NewParser().ParseBytes(content)
	if err != nil {
		return internalReport.NewError(err).
			WithMessage("irregular project manifest")
	}

	// Validate node
	validation, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(manifestSchema),
		internalYaml.NewJsonLoader(man.node),
	)
	if err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to validate project manifest")
	}

	if !validation.Valid() {
		return internalReport.NewError(
			internalValidation.NewError(
				"invalid project manifest",
				validation,
			).
				WithReporter(man).
				WithMessages([]internalValidation.ErrorMessage{
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
				}),
		)
	}

	// Extract config node
	configNode, err := internalYaml.NewExtractor(&man.node).ExtractRootMap("manala")
	if err != nil {
		return internalReport.NewError(err).
			WithMessage("incorrect project manifest")
	}

	// Decode config
	if err = yaml.NodeToValue(configNode, man.config); err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to decode project manifest config")
	}

	// Decode vars
	if err = yaml.NodeToValue(man.node, &man.vars); err != nil {
		return internalReport.NewError(err).
			WithMessage("unable to decode recipe manifest vars")
	}

	return nil
}

func (man *Manifest) Report(result gojsonschema.ResultError, report *internalReport.Report) {
	internalYaml.NewValidationPathReporter(man.node).Report(result, report)
}

type manifestConfig struct {
	Recipe     string
	Repository string
}
