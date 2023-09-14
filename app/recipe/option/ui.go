package option

import (
	"manala/app"
	"manala/internal/serrors"
	"manala/internal/ui/components"
)

func NewUiFormField(option app.RecipeOption, vars *map[string]any) (components.FormField, error) {
	switch option := option.(type) {
	case *SelectOption:
		return NewSelectOptionUiFormField(option, vars)
	case *TextOption:
		return NewTextOptionUiFormField(option, vars)
	}

	return nil, serrors.New("unknown recipe option").
		WithArguments("label", option.Label())
}
