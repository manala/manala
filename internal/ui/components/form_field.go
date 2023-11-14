package components

import (
	"manala/internal/accessor"
	"manala/internal/validator"
)

type FormField interface {
	Name() string
	Label() string
	Help() string
	Set(value any) error
	Submit() (bool, error)
	Violations() validator.Violations
}

func newFormField(name string, label string, help string, accessor accessor.Accessor, validator validator.Validator) (*formField, error) {
	// Initial value
	value, err := accessor.Get()
	if err != nil {
		return nil, err
	}

	return &formField{
		name:      name,
		label:     label,
		help:      help,
		value:     value,
		accessor:  accessor,
		validator: validator,
	}, nil
}

type formField struct {
	name       string
	label      string
	help       string
	value      any
	accessor   accessor.Accessor
	validator  validator.Validator
	violations validator.Violations
}

func (field *formField) Name() string {
	return field.name
}

func (field *formField) Label() string {
	return field.label
}

func (field *formField) Help() string {
	return field.help
}

func (field *formField) Set(value any) error {
	field.value = value
	return field.validate()
}

func (field *formField) Submit() (bool, error) {
	if err := field.validate(); err != nil {
		return false, err
	}
	if len(field.violations) != 0 {
		return false, nil
	}
	if err := field.accessor.Set(field.value); err != nil {
		return false, err
	}
	return len(field.violations) == 0, nil
}

func (field *formField) validate() (err error) {
	field.violations, err = field.validator.Validate(field.value)
	return err
}

func (field *formField) Violations() validator.Violations {
	return field.violations
}
