package option

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/ui/components"
)

func NewUIFormField(option app.RecipeOption, vars *map[string]any) (components.FormField, error) {
	switch option := option.(type) {
	case *SelectOption:
		return NewSelectOptionUIFormField(option, vars)
	case *TextOption:
		return NewTextOptionUIFormField(option, vars)
	}

	return nil, serrors.New("unknown recipe option").
		WithArguments("label", option.Label())
}
