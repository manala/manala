package init

import (
	"github.com/manala/manala/app"

	"codeberg.org/tslocum/cview"
	"github.com/gdamore/tcell/v3/color"
)

type DialogList struct {
	*cview.List
}

func NewDialogList(title string) (*DialogList, *cview.Flex) {
	// List
	list := &DialogList{List: cview.NewList()}
	list.SetPadding(1, 1, 1, 1)
	list.SetHover(true)
	list.SetWrapAround(true)
	list.SetIndicators("▶", "", " ", "")
	list.SetBackgroundColor(color.Default)
	list.SetMainTextColor(DialogStyles.PrimaryColor)
	list.SetSecondaryTextColor(DialogStyles.SecondaryColor)
	list.SetSelectedTextColor(DialogStyles.TertiaryColor)
	list.SetSelectedBackgroundColor(DialogStyles.SecondaryColor)
	list.SetScrollBarColor(DialogStyles.SecondaryColor)

	return list, NewDialogPanel(title, list)
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
