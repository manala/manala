package project

import (
	_ "embed"
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"io"
	"manala/internal/schema"
	"manala/internal/serrors"
	"manala/internal/validator"
	"manala/internal/yaml"
	"regexp"
)

//go:embed resources/manifest.schema.json
var _manifestSchemaSource []byte
var _manifestSchema = schema.MustParse(_manifestSchemaSource)

func NewManifest() *Manifest {
	return &Manifest{
		config: &manifestConfig{},
		vars:   map[string]any{},
	}
}

type Manifest struct {
	node   goYamlAst.Node
	config *manifestConfig
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

func (manifest *Manifest) ReadFrom(reader io.Reader) (n int64, err error) {
	// Read content
	content, err := io.ReadAll(reader)
	n = int64(len(content))
	if err != nil {
		return n, serrors.New("unable to read project manifest").
			WithErrors(err)
	}

	// Parse content to node
	manifest.node, err = yaml.NewParser().ParseBytes(content)
	if err != nil {
		return n, serrors.New("irregular project manifest").
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
			{Path: "manala", Type: validator.REQUIRED, Property: "recipe", StructuredMessage: "missing manala recipe property"},
			{PathRegex: regexp.MustCompile(`^manala\.[^.\[]+$`), Type: validator.ADDITIONAL_PROPERTY_NOT_ALLOWED, StructuredMessage: "manala field don't support additional properties"},
			// Recipe
			{Path: "manala.recipe", Type: validator.INVALID_TYPE, StructuredMessage: "manala recipe field must be a string"},
			{Path: "manala.recipe", Type: validator.STRING_GTE, StructuredMessage: "empty manala recipe field"},
			{Path: "manala.recipe", Type: validator.STRING_LTE, StructuredMessage: "too long manala recipe field"},
			// Repository
			{Path: "manala.repository", Type: validator.INVALID_TYPE, StructuredMessage: "manala repository field must be a string"},
			{Path: "manala.repository", Type: validator.STRING_GTE, StructuredMessage: "empty manala repository field"},
			{Path: "manala.repository", Type: validator.STRING_LTE, StructuredMessage: "too long manala repository field"},
		}),
		validator.WithFormatters(
			manifest.ValidatorFormatter(),
		),
	).Validate(data); err != nil {
		return n, serrors.New("unable to validate project manifest").
			WithErrors(err)
	} else if len(violations) != 0 {
		return n, serrors.New("invalid project manifest").
			WithErrors(violations.StructuredErrors()...)
	}

	// Extract config node
	configNode, err := yaml.NewExtractor(&manifest.node).ExtractRootMap("manala")
	if err != nil {
		return n, serrors.New("incorrect project manifest").
			WithErrors(err)
	}

	// Decode config
	if err = goYaml.NodeToValue(configNode, manifest.config); err != nil {
		return n, serrors.New("unable to decode project manifest config").
			WithErrors(err)
	}

	// Decode vars
	if err = goYaml.NodeToValue(manifest.node, &manifest.vars); err != nil {
		return n, serrors.New("unable to decode recipe manifest vars").
			WithErrors(err)
	}

	return n, err
}

func (manifest *Manifest) ValidatorFormatter() validator.Formatter {
	return yaml.NodeValidatorFormatter(manifest.node)
}

type manifestConfig struct {
	Recipe     string
	Repository string
}
