package manifest

import (
	"errors"

	"github.com/manala/manala/internal/validation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlvalidation "github.com/manala/manala/internal/yaml/validation"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

var configValidator = validation.MustNewValidator(map[string]any{
	"type": "object",
	"properties": map[string]any{
		"recipe":     map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"repository": map[string]any{"type": "string", "minLength": 1, "maxLength": 256},
	},
	"additionalProperties": false,
	"required":             []any{"recipe"},
})

type Config struct {
	Recipe     string `yaml:"recipe"`
	Repository string `yaml:"repository"`
}

func (c *Config) UnmarshalYAML(node ast.Node) error {
	// Decode to map for validation
	var data map[string]any
	if err := yaml.NodeToValue(node, &data); err != nil {
		return yamlerrors.From(err)
	}

	// Validate
	if violations, err := configValidator.Validate(data, yamlvalidation.WithLocator(node)); violations != nil || err != nil {
		return errors.Join(violations, err)
	}

	// Decode using type alias to breaks UnmarshalYAML recursion
	type config Config
	if err := yaml.NodeToValue(node, (*config)(c)); err != nil {
		return yamlerrors.From(err)
	}

	return nil
}
