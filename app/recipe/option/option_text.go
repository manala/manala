package option

import (
	"github.com/manala/manala/internal/json"
	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/ui/components"
	"github.com/manala/manala/internal/validator"
)

func NewTextOptionUIFormField(option *TextOption, vars *map[string]any) (components.FormField, error) {
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

func NewTextOption(option *option, _ map[string]any) (*TextOption, error) {
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
