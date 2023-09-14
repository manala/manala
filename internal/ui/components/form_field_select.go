package components

import (
	"fmt"
	"manala/internal/accessor"
	"manala/internal/serrors"
	"manala/internal/validator"
	"slices"
)

func NewFormFieldSelect(name string, label string, help string, accessor accessor.Accessor, validator validator.Validator) (*FormFieldSelect, error) {
	field, err := newFormField(name, label, help, accessor, validator)
	if err != nil {
		return nil, err
	}

	return &FormFieldSelect{
		formField: field,
	}, nil
}

type FormFieldSelect struct {
	Options []*FormFieldSelectOption
	*formField
}

func (field *FormFieldSelect) GetIndex() int {
	return slices.IndexFunc(field.Options, func(option *FormFieldSelectOption) bool {
		return option.Value() == field.value
	})
}

func (field *FormFieldSelect) SetIndex(index int) error {
	if index < 0 || index >= len(field.Options) {
		return serrors.New("invalid select index")
	}
	field.value = field.Options[index].Value()
	return nil
}

/**********/
/* Option */
/**********/

func NewFormFieldSelectOption(value any) *FormFieldSelectOption {
	return &FormFieldSelectOption{
		value: value,
	}
}

type FormFieldSelectOption struct {
	value any
}

func (option *FormFieldSelectOption) Value() any {
	return option.value
}

func (option *FormFieldSelectOption) Label() string {
	switch option.value {
	case nil:
		return "<None>"
	case true:
		return "<True>"
	case false:
		return "<False>"
	}

	return fmt.Sprintf("%v", option.value)
}
