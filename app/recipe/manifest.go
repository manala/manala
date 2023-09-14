package recipe

import (
	_ "embed"
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"io"
	"manala/app"
	"manala/app/recipe/option"
	"manala/internal/schema"
	"manala/internal/serrors"
	"manala/internal/syncer"
	"manala/internal/validator"
	"manala/internal/yaml"
	"regexp"
)

//go:embed resources/manifest.schema.json
var _manifestSchemaSource []byte
var _manifestSchema = schema.MustParse(_manifestSchemaSource)

func NewManifest() *Manifest {
	return &Manifest{
		config: &manifestConfig{
			Sync: sync{},
		},
		vars:   map[string]any{},
		schema: schema.Schema{},
	}
}

type Manifest struct {
	node    goYamlAst.Node
	config  *manifestConfig
	vars    map[string]any
	schema  schema.Schema
	options []app.RecipeOption
}

func (manifest *Manifest) Description() string {
	return manifest.config.Description
}

func (manifest *Manifest) Template() string {
	return manifest.config.Template
}

func (manifest *Manifest) Vars() map[string]any {
	return manifest.vars
}

func (manifest *Manifest) Sync() []syncer.UnitInterface {
	return manifest.config.Sync
}

func (manifest *Manifest) Schema() schema.Schema {
	return manifest.schema
}

func (manifest *Manifest) Options() []app.RecipeOption {
	return manifest.options
}

func (manifest *Manifest) ReadFrom(reader io.Reader) (n int64, err error) {
	// Read content
	content, err := io.ReadAll(reader)
	n = int64(len(content))
	if err != nil {
		return n, serrors.New("unable to read recipe manifest").
			WithErrors(err)
	}

	// Parse content to node
	manifest.node, err = yaml.NewParser(yaml.WithComments()).ParseBytes(content)
	if err != nil {
		return n, serrors.New("irregular recipe manifest").
			WithErrors(err)
	}

	// Decode node
	var data any
	if err := goYaml.NewDecoder(manifest.node).Decode(&data); err != nil {
		// Nil or empty content
		if err == io.EOF {
			return n, serrors.New("empty content")
		}
		return n, yaml.NewError(err)
	}

	// Validate node data
	if violations, err := validator.New(
		validator.WithValidators(
			schema.NewValidator(_manifestSchema),
		),
		validator.WithFilters(validator.Filters{
			{Path: "", Type: validator.INVALID_TYPE, StructuredMessage: "yaml document must be a map"},
			{Path: "", Type: validator.REQUIRED, Property: "manala", StructuredMessage: "missing manala property"},
			{Path: "manala", Type: validator.INVALID_TYPE, StructuredMessage: "manala field must be a map"},
			{Path: "manala", Type: validator.REQUIRED, Property: "description", StructuredMessage: "missing manala description property"},
			{PathRegex: regexp.MustCompile(`^manala\.[^.\[]+$`), Type: validator.ADDITIONAL_PROPERTY_NOT_ALLOWED, StructuredMessage: "manala field don't support additional properties"},
			// Description
			{Path: "manala.description", Type: validator.INVALID_TYPE, StructuredMessage: "manala description field must be a string"},
			{Path: "manala.description", Type: validator.STRING_GTE, StructuredMessage: "empty manala description field"},
			{Path: "manala.description", Type: validator.STRING_LTE, StructuredMessage: "too long manala description field"},
			// Template
			{Path: "manala.template", Type: validator.INVALID_TYPE, StructuredMessage: "manala template field must be a string"},
			{Path: "manala.template", Type: validator.STRING_GTE, StructuredMessage: "empty manala template field"},
			{Path: "manala.template", Type: validator.STRING_LTE, StructuredMessage: "too long manala template field"},
			// Sync
			{Path: "manala.sync", Type: validator.INVALID_TYPE, StructuredMessage: "manala sync field must be a sequence"},
			// Sync Item
			{PathRegex: regexp.MustCompile(`^manala\.sync\[\d+]$`), Type: validator.INVALID_TYPE, StructuredMessage: "manala sync sequence entries must be strings"},
			{PathRegex: regexp.MustCompile(`^manala\.sync\[\d+]$`), Type: validator.STRING_GTE, StructuredMessage: "empty manala sync sequence entry"},
			{PathRegex: regexp.MustCompile(`^manala\.sync\[\d+]$`), Type: validator.STRING_LTE, StructuredMessage: "too long manala sync sequence entry"},
		}),
		validator.WithFormatters(
			manifest.ValidatorFormatter(),
		),
	).Validate(data); err != nil {
		return n, serrors.New("unable to validate recipe manifest").
			WithErrors(err)
	} else if len(violations) != 0 {
		return n, serrors.New("invalid recipe manifest").
			WithErrors(violations.StructuredErrors()...)
	}

	// Extract config node
	configNode, err := yaml.NewExtractor(&manifest.node).ExtractRootMap("manala")
	if err != nil {
		return n, serrors.New("incorrect recipe manifest").
			WithErrors(err)
	}

	// Decode config
	if err = goYaml.NodeToValue(configNode, manifest.config); err != nil {
		return n, serrors.New("unable to decode recipe manifest config").
			WithErrors(err)
	}

	// Decode vars
	if err = goYaml.NodeToValue(manifest.node, &manifest.vars); err != nil {
		return n, serrors.New("unable to decode recipe manifest vars").
			WithErrors(err)
	}

	// Infer schema
	if err = yaml.NewNodeSchemaInferrer(manifest.node).Infer(manifest.schema); err != nil {
		return n, serrors.New("unable to infer recipe manifest schema").
			WithErrors(err)
	}

	// Infer options
	if err = option.NewInferrer().Infer(manifest.node, &manifest.options); err != nil {
		return n, serrors.New("unable to infer recipe manifest options").
			WithErrors(err)
	}

	return n, err
}

func (manifest *Manifest) ValidatorFormatter() validator.Formatter {
	return yaml.NodeValidatorFormatter(manifest.node)
}

type manifestConfig struct {
	Description string
	Template    string
	Sync        sync
}
