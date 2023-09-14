package option

import (
	"manala/internal/json"
	"manala/internal/path"
	"manala/internal/schema"
	"manala/internal/serrors"
	"manala/internal/ui/components"
	"manala/internal/validator"
)

func NewSelectOption(option *option, fields map[string]any) (*SelectOption, error) {
	// Option
	selectOption := &SelectOption{
		option: option,
	}

	// Enum
	enum, ok := option.schema["enum"].([]any)
	if !ok {
		return nil, serrors.New("invalid recipe option enum").
			WithArguments("label", option.label)
	}
	if len(enum) == 0 {
		return nil, serrors.New("empty recipe option enum").
			WithArguments("label", option.label)
	}

	// Values
	selectOption.Values = make([]any, len(enum))
	for i := range enum {
		if value, ok := json.NumberType(enum[i]); ok {
			selectOption.Values[i] = value.Normalize()
		} else {
			selectOption.Values[i] = enum[i]
		}
	}

	return selectOption, nil
}

func NewSelectOptionUiFormField(option *SelectOption, vars *map[string]any) (components.FormField, error) {
	// Field
	field, err := components.NewFormFieldSelect(
		option.Name(),
		option.Label(),
		option.Help(),
		path.NewAccessor(
			option.Path(),
			vars,
		),
		validator.New(
			validator.WithValidators(
				schema.NewValidator(option.Schema()),
				option,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// Options
	field.Options = make([]*components.FormFieldSelectOption, len(option.Values))
	for i := range option.Values {
		field.Options[i] = components.NewFormFieldSelectOption(option.Values[i])
	}

	return field, nil
}

type SelectOption struct {
	*option
	Values []any
}
