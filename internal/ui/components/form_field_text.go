package components

import (
	"github.com/manala/manala/internal/accessor"
	"github.com/manala/manala/internal/validator"
)

func NewFormFieldText(name string, label string, help string, accessor accessor.Accessor, validator validator.Validator) (*FormFieldText, error) {
	field, err := newFormField(name, label, help, accessor, validator)
	if err != nil {
		return nil, err
	}

	return &FormFieldText{
		formField: field,
	}, nil
}

type FormFieldText struct {
	MaxLength int
	*formField
}

func (field *FormFieldText) Get() any {
	return field.value
}
