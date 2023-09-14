package api

import (
	"manala/app"
	"manala/app/recipe/option"
	"manala/internal/accessor"
	"manala/internal/serrors"
	"manala/internal/ui/components"
)

func (api *Api) LoadPrecedingRecipe(repository app.Repository) (app.Recipe, error) {
	// Log
	api.log.Debug("load preceding recipe…")

	// Load preceding recipe
	return api.recipeManager.LoadPrecedingRecipe(repository)
}

// NewUiRecipeList create an ui recipe list
func (api *Api) NewUiRecipeList(items []app.Recipe) []components.ListItem {
	list := make([]components.ListItem, len(items))
	for i, item := range items {
		list[i] = components.ListItem{
			Primary:   item.Name(),
			Secondary: item.Description(),
		}
	}

	return list
}

// NewUiRecipeListForm create an ui recipe list form
func (api *Api) NewUiRecipeListForm(items []app.Recipe, value *app.Recipe) (*components.ListForm, error) {
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

// NewUiRecipeOptionsForm create an ui recipe options form
func (api *Api) NewUiRecipeOptionsForm(recipe app.Recipe, vars *map[string]any) (*components.Form, error) {
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
