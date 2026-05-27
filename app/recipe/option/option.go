package option

import (
	"github.com/manala/manala/internal/validation"
)

var Validator = validation.MustNewValidator(map[string]any{
	"type": "object",
	"properties": map[string]any{
		"type":  map[string]any{"enum": []any{STRING, ENUM}},
		"name":  map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"label": map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"help":  map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
	},
	"required":             []any{"label"},
	"additionalProperties": false,
})

type option struct {
	Name  string `yaml:"name"`
	Label string `yaml:"label"`
	Help  string `yaml:"help"`
}
