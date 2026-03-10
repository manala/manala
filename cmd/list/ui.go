package list

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/ui/components"
)

func NewUIRecipeList(items []app.Recipe) []components.ListItem {
	list := make([]components.ListItem, len(items))
	for i, item := range items {
		list[i] = components.ListItem{
			Primary:   item.Name(),
			Secondary: item.Description(),
		}
	}

	return list
}
