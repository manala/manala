package manifest

import (
	_ "embed"
	"io"
	"regexp"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/sync"
	"github.com/manala/manala/internal/validator"
	"github.com/manala/manala/internal/yaml"
	"github.com/manala/manala/internal/yaml/parser"

	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
)

const filename = ".manala.yaml"

//go:embed schema.json
var _schemaSource []byte
var _schema = schema.MustParse(_schemaSource)

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
		return serrors.New("irregular recipe manifest").
			WithErrors(err)
	}

	// Decode node
	var data any
	if err := goYaml.NewDecoder(manifest.node).Decode(&data); err != nil {
		// Nil or empty content
		if err == io.EOF {
			return serrors.New("empty content")
		}

		return yaml.NewError(err)
	}

	// Validate node data
	if violations, err := validator.New(
		validator.WithValidators(
			schema.NewValidator(_schema),
		),
		validator.WithFilters(validator.Filters{
			{Path: "", Type: validator.Required, Property: "manala", StructuredMessage: "missing manala property"},
			{Path: "manala", Type: validator.InvalidType, StructuredMessage: "manala field must be a map"},
			{Path: "manala", Type: validator.Required, Property: "description", StructuredMessage: "missing manala description property"},
			{PathRegex: regexp.MustCompile(`^manala\.[^.\[]+$`), Type: validator.AdditionalPropertyNotAllowed, StructuredMessage: "manala field don't support additional properties"},
			// Description
			{Path: "manala.description", Type: validator.InvalidType, StructuredMessage: "manala description field must be a string"},
			{Path: "manala.description", Type: validator.StringGte, StructuredMessage: "empty manala description field"},
			{Path: "manala.description", Type: validator.StringLte, StructuredMessage: "too long manala description field"},
			// Icon
			{Path: "manala.icon", Type: validator.InvalidType, StructuredMessage: "manala icon field must be a string"},
			{Path: "manala.icon", Type: validator.StringGte, StructuredMessage: "empty manala icon field"},
			{Path: "manala.icon", Type: validator.StringLte, StructuredMessage: "too long manala icon field"},
			// Template
			{Path: "manala.template", Type: validator.InvalidType, StructuredMessage: "manala template field must be a string"},
			{Path: "manala.template", Type: validator.StringGte, StructuredMessage: "empty manala template field"},
			{Path: "manala.template", Type: validator.StringLte, StructuredMessage: "too long manala template field"},
			// Sync
			{Path: "manala.sync", Type: validator.InvalidType, StructuredMessage: "manala sync field must be a sequence"},
			// Sync Item
			{PathRegex: regexp.MustCompile(`^manala\.sync\[\d+]$`), Type: validator.InvalidType, StructuredMessage: "manala sync sequence entries must be strings"},
			{PathRegex: regexp.MustCompile(`^manala\.sync\[\d+]$`), Type: validator.StringGte, StructuredMessage: "empty manala sync sequence entry"},
			{PathRegex: regexp.MustCompile(`^manala\.sync\[\d+]$`), Type: validator.StringLte, StructuredMessage: "too long manala sync sequence entry"},
		}),
		validator.WithFormatters(
			manifest.ValidatorFormatter(),
		),
	).Validate(data); err != nil {
		return serrors.New("unable to validate recipe manifest").
			WithErrors(err)
	} else if len(violations) != 0 {
		return serrors.New("invalid recipe manifest").
			WithErrors(violations.StructuredErrors()...)
	}

	// Extract config node
	configNode, err := yaml.NewExtractor(manifest.node).ExtractRootMap("manala")
	if err != nil {
		return serrors.New("incorrect recipe manifest").
			WithErrors(err)
	}

	// Decode config
	if err = goYaml.NodeToValue(configNode, manifest.config); err != nil {
		return serrors.New("unable to decode recipe manifest config").
			WithErrors(err)
	}

	// Decode vars
	if err = goYaml.NodeToValue(manifest.node, &manifest.vars); err != nil {
		return serrors.New("unable to decode recipe manifest vars").
			WithErrors(err)
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

func (manifest *Manifest) ValidatorFormatter() validator.Formatter {
	return yaml.NodeValidatorFormatter(manifest.node)
}

type config struct {
	Description string
	Icon        string
	Template    string
	Sync        recipe.Sync
}
