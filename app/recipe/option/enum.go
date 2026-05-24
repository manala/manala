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

const ENUM = "enum"

var enumValidator = validation.MustNewValidator(map[string]any{
	"type": "object",
	"properties": map[string]any{
		"type":  map[string]any{"const": "enum"},
		"name":  map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"label": map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
		"help":  map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
	},
	"additionalProperties": false,
	"required":             []any{"label"},
})

type Enum struct {
	name    string
	label   string
	help    string
	values  []any
	pointer jsonpointer.Pointer
}

func NewEnum(sch map[string]any, path string) (*Enum, error) {
	// Schema *MUST* contains enum
	enum, ok := sch["enum"].([]any)
	if !ok {
		return nil, errors.New("invalid recipe option enum")
	}

	if len(enum) == 0 {
		return nil, errors.New("empty recipe option enum")
	}

	o := &Enum{}

	// Values
	o.values = make([]any, len(enum))
	for i := range enum {
		if value, ok := jsonnumber.NumberType(enum[i]); ok {
			o.values[i] = value.Normalize()
		} else {
			o.values[i] = enum[i]
		}
	}

	// Pointer
	var err error
	if o.pointer, err = jsonpointer.New(yamlpath.ToJSONPointer(path)); err != nil {
		return nil, err
	}

	return o, nil
}

func (o *Enum) Name() string  { return o.name }
func (o *Enum) Label() string { return o.label }
func (o *Enum) Help() string  { return o.help }
func (o *Enum) Values() []any { return o.values }

func (o *Enum) Get(data *map[string]any) (any, error) {
	value, _, err := o.pointer.Get(data)
	return value, err
}

func (o *Enum) Set(data *map[string]any, v any) error {
	_, err := o.pointer.Set(data, v)
	return err
}

func (o *Enum) UnmarshalJSON(bytes []byte) error {
	// Decode to map for validation
	var data map[string]any
	if err := jsondecoder.Decode(bytes, &data); err != nil {
		return err
	}

	// Validate
	if err := enumValidator.Validate(data, jsonvalidation.WithLocator(bytes)); err != nil {
		return err
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
