package option

import (
	"errors"

	jsondecoder "github.com/manala/manala/internal/json/decoder"
	jsonnumber "github.com/manala/manala/internal/json/number"
	jsonvalidation "github.com/manala/manala/internal/json/validation"
	"github.com/manala/manala/internal/validation"
	yamlpath "github.com/manala/manala/internal/yaml/path"

	"github.com/go-openapi/jsonpointer"
	"github.com/gosimple/slug"
)

const STRING = "string"

var stringValidator = validation.MustNewValidator(map[string]any{
	"type": "object",
	"properties": map[string]any{
		"type":  map[string]any{"const": "string"},
		"name":  map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"label": map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"help":  map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
	},
	"additionalProperties": false,
	"required":             []any{"label"},
})

type String struct {
	name      string
	label     string
	help      string
	maxLength int
	pointer   jsonpointer.Pointer
	validator *validation.Validator
}

func NewString(sch map[string]any, path string) (*String, error) {
	// Schema type *MUST* be string
	if t, ok := sch["type"]; !ok || t != "string" {
		return nil, errors.New("invalid recipe option string type")
	}

	o := &String{}

	// Max length
	if maxLength, ok := jsonnumber.NumberType(sch["maxLength"]); ok {
		o.maxLength = maxLength.Int()
	}

	// Pointer
	var err error
	if o.pointer, err = jsonpointer.New(yamlpath.ToJSONPointer(path)); err != nil {
		return nil, err
	}

	// Validator
	if o.validator, err = validation.NewValidator(sch); err != nil {
		return nil, err
	}

	return o, nil
}

func (o *String) Name() string   { return o.name }
func (o *String) Label() string  { return o.label }
func (o *String) Help() string   { return o.help }
func (o *String) MaxLength() int { return o.maxLength }

func (o *String) Get(data *map[string]any) (string, error) {
	value, _, err := o.pointer.Get(data)
	if err != nil {
		return "", err
	}
	if value, ok := value.(string); ok {
		return value, nil
	}
	return "", nil
}

func (o *String) Set(data *map[string]any, v string) error {
	_, err := o.pointer.Set(data, v)
	return err
}

func (o *String) Validate(v string) error {
	violations, err := o.validator.Validate(v)
	if err != nil {
		return err
	}
	if violation, ok := violations.First(); ok {
		return violation
	}
	return nil
}

func (o *String) UnmarshalJSON(bytes []byte) error {
	// Decode to map for validation
	var data map[string]any
	if err := jsondecoder.Decode(bytes, &data); err != nil {
		return err
	}

	// Validate
	if violations, err := stringValidator.Validate(data, jsonvalidation.WithLocator(bytes)); violations != nil || err != nil {
		return errors.Join(violations, err)
	}

	// Decode
	var env struct {
		Name  string `json:"name"`
		Label string `json:"label"`
		Help  string `json:"help"`
	}
	if err := jsondecoder.Decode(bytes, &env); err != nil {
		return err
	}

	o.label = env.Label
	o.help = env.Help
	if o.name = env.Name; o.name == "" {
		o.name = slug.Make(o.label)
	}

	return nil
}
