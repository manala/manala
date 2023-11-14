package option

import (
	"manala/internal/json"
	"manala/internal/path"
	"manala/internal/schema"
	"manala/internal/serrors"
	"manala/internal/ui/components"
	"manala/internal/validator"
)

func NewTextOption(option *option, fields map[string]any) (*TextOption, error) {
	// Schema type *MUST* be string
	if t, ok := option.schema["type"]; !ok || t != "string" {
		return nil, serrors.New("invalid recipe option string type").
			WithArguments("label", option.label)
	}

	// Option
	textOption := &TextOption{
		option: option,
	}

	// Max length
	if maxLength, ok := json.NumberType(option.schema["maxLength"]); ok {
		textOption.MaxLength = maxLength.Int()
	}

	return textOption, nil
}

func NewTextOptionUiFormField(option *TextOption, vars *map[string]any) (components.FormField, error) {
	// Field
	field, err := components.NewFormFieldText(
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

	// Max length
	field.MaxLength = option.MaxLength

	return field, nil
}

type TextOption struct {
	*option
	MaxLength int
}
