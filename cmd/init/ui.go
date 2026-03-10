package init

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/accessor"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/ui/components"
)

func NewUIRecipeListForm(items []app.Recipe, value *app.Recipe) (*components.ListForm, error) {
	// List
	list := make([]components.ListItem, len(items))
	for i, item := range items {
		list[i] = components.ListItem{
			Primary:   item.Name(),
			Secondary: item.Description(),
		}
	}

	return components.NewListForm(
		list,
		accessor.NewIndex(value, items),
	)
}

func NewUIRecipeOptionsForm(recipe app.Recipe, vars *map[string]any) (*components.Form, error) {
	options := recipe.Options()
	fields := make([]components.FormField, len(options))

	for i := range options {
		var err error
		if fields[i], err = option.NewUIFormField(options[i], vars); err != nil {
			return nil, serrors.New("unable to get recipe option form field").
				WithArguments("label", options[i].Label()).
				WithErrors(err)
		}
	}

	return components.NewForm(fields), nil
}
