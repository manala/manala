package option

import (
	_ "embed"
	"github.com/gosimple/slug"
	"io"
	"manala/app"
	"manala/internal/json"
	"manala/internal/path"
	"manala/internal/schema"
	"manala/internal/serrors"
	"manala/internal/validator"
)

//go:embed resources/schema.json
var _schemaSource []byte
var _schema = schema.MustParse(_schemaSource)

type option struct {
	name   string
	label  string
	help   string
	schema schema.Schema
	path   path.Path
}

func (option *option) Name() string {
	return option.name
}

func (option *option) Label() string {
	return option.label
}

func (option *option) Help() string {
	return option.help
}

func (option *option) Path() path.Path {
	return option.path
}

func (option *option) Schema() schema.Schema {
	return option.schema
}

func (option *option) Validate(_ any) (validator.Violations, error) {
	return nil, nil
}

func New(reader io.Reader, optionSchema schema.Schema, optionPath path.Path) (app.RecipeOption, error) {
	// Read content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, serrors.New("unable to read recipe option").
			WithErrors(err)
	}

	// Parse content to fields
	var fields map[string]any
	if err := json.Unmarshal(content, &fields); err != nil {
		return nil, serrors.New("irregular recipe option").
			WithErrors(err)
	}

	// Validate fields
	if violations, err := validator.New(
		validator.WithValidators(
			schema.NewValidator(_schema),
		),
	).Validate(fields); err != nil {
		return nil, serrors.New("unable to validate recipe option").
			WithErrors(err)
	} else if len(violations) != 0 {
		return nil, serrors.New("invalid recipe option").
			WithErrors(violations.StructuredErrors()...)
	}

	option := &option{
		label:  fields["label"].(string),
		schema: optionSchema,
		path:   optionPath,
	}

	// Name
	var ok bool
	if option.name, ok = fields["name"].(string); !ok {
		option.name = slug.Make(option.label)
	}

	// Help
	option.help, _ = fields["help"].(string)

	// Type
	var optionType string
	if optionType, ok = fields["type"].(string); !ok {
		// Auto detection
		if _, ok := optionSchema["enum"]; ok {
			optionType = "select"
		} else if schemaType, ok := optionSchema["type"]; ok && schemaType == "string" {
			optionType = "text"
		} else {
			return nil, serrors.New("unable to auto detect recipe option type").
				WithArguments("label", option.label)
		}
	}

	switch optionType {
	case "text":
		return NewTextOption(option, fields)
	case "select":
		return NewSelectOption(option, fields)
	}

	return nil, serrors.New("unknown recipe option type").
		WithArguments("label", option.label)
}
