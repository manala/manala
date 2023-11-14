package charm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"manala/internal/ui/components"
)

func (adapter *Adapter) List(header string, list []components.ListItem) error {
	renderer := adapter.outRenderer

	// Items views
	views := make([]string, len(list))
	for i, item := range list {
		views[i] = newListItemModel(
			item,
			listItemStyle.New(renderer),
			listItemPrimaryStyle.New(renderer),
			listItemSecondaryStyle.New(renderer),
		).View()
	}

	_, _ = renderer.Output().WriteString(
		lipgloss.JoinVertical(lipgloss.Left,
			// Header
			newHeaderModel(header, renderer).View(),
			// Views
			lipgloss.JoinVertical(lipgloss.Left,
				views...,
			),
		) + "\n",
	)

	return nil
}

/********/
/* Item */
/********/

func newListItemModel(item components.ListItem, style *style, primaryStyle *style, secondaryStyle *style) listItemModel {
	return listItemModel{
		item:           item,
		style:          style,
		primaryStyle:   primaryStyle,
		secondaryStyle: secondaryStyle,
	}
}

type listItemModel struct {
	item           components.ListItem
	style          *style
	primaryStyle   *style
	secondaryStyle *style
}

func (model listItemModel) Init() tea.Cmd {
	return nil
}

func (model listItemModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update styles
	model.style.Update(msg)
	model.primaryStyle.Update(msg)
	model.secondaryStyle.Update(msg)

	return model, nil
}

func (model listItemModel) View() string {
	// Primary
	view := model.primaryStyle.Render(
		model.item.Primary,
	)

	// Secondary
	secondary := model.item.Secondary
	if secondary != "" {
		view = lipgloss.JoinVertical(lipgloss.Left,
			view,
			model.secondaryStyle.Render(
				secondary,
			),
		)
	}

	return model.style.Render(view)
}
