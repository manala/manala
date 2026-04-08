package manifest

import (
	_ "embed"
	"io"
	"regexp"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/validator"
	"github.com/manala/manala/internal/yaml"
	"github.com/manala/manala/internal/yaml/parser"

	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
)

const filename = ".manala.yaml"

//go:embed template.yaml.tmpl
var _template string

//go:embed schema.json
var _schemaSource []byte
var _schema = schema.MustParse(_schemaSource)

func New() *Manifest {
	return &Manifest{
		config: &config{},
		vars:   map[string]any{},
	}
}

type Manifest struct {
	node   *goYamlAst.MappingNode
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

	// Decode node
	var data any
	if err := goYaml.NewDecoder(manifest.node).Decode(&data); err != nil {
		// Nil or empty content
		if err == io.EOF {
			return &parsing.Error{
				Err: serrors.New("empty yaml content"),
			}
		}

		return parser.ErrorFrom(err)
	}

	// Validate node data
	if violations, err := validator.New(
		validator.WithValidators(
			schema.NewValidator(_schema),
		),
		validator.WithFilters(validator.Filters{
			{Path: "", Type: validator.Required, Property: "manala", StructuredMessage: "missing manala property"},
			{Path: "manala", Type: validator.InvalidType, StructuredMessage: "manala field must be a map"},
			{Path: "manala", Type: validator.Required, Property: "recipe", StructuredMessage: "missing manala recipe property"},
			{PathRegex: regexp.MustCompile(`^manala\.[^.\[]+$`), Type: validator.AdditionalPropertyNotAllowed, StructuredMessage: "manala field don't support additional properties"},
			// Recipe
			{Path: "manala.recipe", Type: validator.InvalidType, StructuredMessage: "manala recipe field must be a string"},
			{Path: "manala.recipe", Type: validator.StringGte, StructuredMessage: "empty manala recipe field"},
			{Path: "manala.recipe", Type: validator.StringLte, StructuredMessage: "too long manala recipe field"},
			// Repository
			{Path: "manala.repository", Type: validator.InvalidType, StructuredMessage: "manala repository field must be a string"},
			{Path: "manala.repository", Type: validator.StringGte, StructuredMessage: "empty manala repository field"},
			{Path: "manala.repository", Type: validator.StringLte, StructuredMessage: "too long manala repository field"},
		}),
		validator.WithFormatters(
			manifest.ValidatorFormatter(),
		),
	).Validate(data); err != nil {
		return serrors.New("unable to validate project manifest").
			WithErrors(err)
	} else if len(violations) != 0 {
		return serrors.New("invalid project manifest").
			WithErrors(violations.StructuredErrors()...)
	}

	// Extract config node
	configNode, err := yaml.NewExtractor(manifest.node).ExtractRootMap("manala")
	if err != nil {
		return serrors.New("incorrect project manifest").
			WithErrors(err)
	}

	// Decode config
	if err = goYaml.NodeToValue(configNode, manifest.config); err != nil {
		return serrors.New("unable to decode project manifest config").
			WithErrors(err)
	}

	// Decode vars
	if err = goYaml.NodeToValue(manifest.node, &manifest.vars); err != nil {
		return serrors.New("unable to decode recipe manifest vars").
			WithErrors(err)
	}

	return nil
}

func (manifest *Manifest) ValidatorFormatter() validator.Formatter {
	return yaml.NodeValidatorFormatter(manifest.node)
}

type config struct {
	Recipe     string
	Repository string
}
