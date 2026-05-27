package option

import (
	"errors"

	jsondecoder "github.com/manala/manala/internal/json/decoder"
	jsonnumber "github.com/manala/manala/internal/json/number"
	"github.com/manala/manala/internal/validation"
	yamlpath "github.com/manala/manala/internal/yaml/path"

	"github.com/go-openapi/jsonpointer"
	"github.com/gosimple/slug"
)

const STRING = "string"

type String struct {
	option    option
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

func (o *String) UnmarshalJSON(bytes []byte) error {
	// Decode as generic option
	if err := jsondecoder.Decode(bytes, &o.option); err != nil {
		return err
	}

	return nil
}

func (o *String) Name() string {
	if o.option.Name == "" {
		o.option.Name = slug.Make(o.option.Label)
	}
	return o.option.Name
}

func (o *String) Label() string { return o.option.Label }
func (o *String) Help() string  { return o.option.Help }

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
	if err := o.validator.Validate(v); err != nil {
		if violations, ok := errors.AsType[validation.Violations](err); ok {
			if violation, ok := violations.First(); ok {
				return violation
			}
		}
		return err
	}
	return nil
}
