package init

import (
	"manala/app"
	"manala/app/recipe/option"
	"manala/internal/accessor"
	"manala/internal/serrors"
	"manala/internal/ui/components"
)

func NewUiRecipeListForm(items []app.Recipe, value *app.Recipe) (*components.ListForm, error) {
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

func NewUiRecipeOptionsForm(recipe app.Recipe, vars *map[string]any) (*components.Form, error) {
	options := recipe.Options()
	fields := make([]components.FormField, len(options))

	for i := range options {
		var err error
		if fields[i], err = option.NewUiFormField(options[i], vars); err != nil {
			return nil, serrors.New("unable to get recipe option form field").
				WithArguments("label", options[i].Label()).
				WithErrors(err)
		}
	}

	return components.NewForm(fields), nil
}
