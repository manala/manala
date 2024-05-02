package list

import (
	"manala/app"
	"manala/internal/ui/components"
)

func NewUiRecipeList(items []app.Recipe) []components.ListItem {
	list := make([]components.ListItem, len(items))
	for i, item := range items {
		list[i] = components.ListItem{
			Primary:   item.Name(),
			Secondary: item.Description(),
		}
	}

	return list
}
