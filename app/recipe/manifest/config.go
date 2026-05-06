package manifest

import (
	"errors"

	"github.com/manala/manala/app/sync"
	"github.com/manala/manala/internal/validation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlvalidation "github.com/manala/manala/internal/yaml/validation"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

var configValidator = validation.MustNewValidator(map[string]any{
	"type": "object",
	"properties": map[string]any{
		"description": map[string]any{"type": "string", "minLength": 1, "maxLength": 256},
		"icon":        map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"template":    map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"partials": map[string]any{
			"type":  "array",
			"items": map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		},
		"sync": map[string]any{},
	},
	"additionalProperties": false,
	"required":             []any{"description"},
})

type Config struct {
	Description string      `yaml:"description"`
	Icon        string      `yaml:"icon"`
	Template    string      `yaml:"template"`
	Partials    []string    `yaml:"partials"`
	Sync        []sync.Unit `yaml:"sync"`
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
