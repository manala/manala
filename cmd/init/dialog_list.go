package init

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/output"

	"codeberg.org/tslocum/cview"
	"github.com/gdamore/tcell/v3/color"
)

type DialogList struct {
	*cview.List
}

func NewDialogList(title string, profile output.Profile) (*DialogList, *cview.Flex) {
	// List
	list := &DialogList{List: cview.NewList()}
	list.SetPadding(1, 1, 1, 1)
	list.SetHover(true)
	list.SetWrapAround(true)
	list.SetIndicators("▶", "", " ", "")
	list.SetBackgroundColor(color.Default)
	list.SetMainTextColor(profile.Color())
	list.SetSecondaryTextColor(profile.MutedColor())
	list.SetSelectedTextColor(profile.ReverseColor())
	list.SetSelectedBackgroundColor(profile.MutedColor())
	list.SetScrollBarColor(profile.MutedColor())

	return list, NewDialogPanel(title, list, profile)
}

func (list *DialogList) SetSelectedFunc(handler func(app.Recipe)) {
	list.List.SetSelectedFunc(func(_ int, item *cview.ListItem) {
		recipe := item.GetReference().(app.Recipe)
		handler(recipe)
	})
}

func (list *DialogList) Build(recipes []app.Recipe) {
	for _, recipe := range recipes {
		item := cview.NewListItem(recipe.Name())
		item.SetSecondaryText("   " + recipe.Description())
		item.SetReference(recipe)

		list.AddItem(item)
	}
}
