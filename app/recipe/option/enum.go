package option

import (
	"errors"

	jsondecoder "github.com/manala/manala/internal/json/decoder"
	jsonnumber "github.com/manala/manala/internal/json/number"
	yamlpath "github.com/manala/manala/internal/yaml/path"

	"github.com/go-openapi/jsonpointer"
	"github.com/gosimple/slug"
)

const ENUM = "enum"

type Enum struct {
	option  option
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

func (o *Enum) UnmarshalJSON(bytes []byte) error {
	// Decode as generic option
	if err := jsondecoder.Decode(bytes, &o.option); err != nil {
		return err
	}

	return nil
}

func (o *Enum) Name() string {
	if o.option.Name == "" {
		o.option.Name = slug.Make(o.option.Label)
	}
	return o.option.Name
}

func (o *Enum) Label() string { return o.option.Label }
func (o *Enum) Help() string  { return o.option.Help }

func (o *Enum) Values() []any { return o.values }

func (o *Enum) Get(data *map[string]any) (any, error) {
	value, _, err := o.pointer.Get(data)
	return value, err
}

func (o *Enum) Set(data *map[string]any, v any) error {
	_, err := o.pointer.Set(data, v)
	return err
}
