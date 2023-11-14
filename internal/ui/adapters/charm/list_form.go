package charm

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"manala/internal/ui"
	"manala/internal/ui/components"
)

func (adapter *Adapter) ListForm(header string, form *components.ListForm) error {
	renderer := adapter.outRenderer

	zone := zone.New()
	defer zone.Close()

	model := listFormModel{
		form:  form,
		items: make([]tea.Model, len(form.List)),
		zone:  zone,
	}

	// Build items
	for i, item := range form.List {
		model.items[i] = newListItemModel(item,
			listFormItemStyle.New(renderer),
			listFormItemPrimaryStyle.New(renderer),
			listFormItemSecondaryStyle.New(renderer),
		)
	}

	// Index
	model.index = newModelsIndex(&model.items)
	model.index.Circular = true

	_model, err := tea.NewProgram(
		newWindowModel(header, model, renderer),
		tea.WithInput(adapter.in),
		tea.WithOutput(renderer.Output()),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
		tea.WithoutSignalHandler(),
	).Run()

	if err != nil {
		return err
	}

	return _model.(windowModel).err
}

/*********/
/* Model */
/*********/

type listFormModel struct {
	form  *components.ListForm
	items []tea.Model
	index *modelsIndex
	zone  *zone.Manager
}

func (model listFormModel) Init() tea.Cmd {
	cmds := newCmds().
		Init(model.items...)

	// Initial index
	index := model.form.GetIndex()
	if index == -1 {
		cmds.Add(
			model.updateIndex(0),
		)
	} else {
		model.setIndex(index)
	}

	return cmds.Batch()
}

func (model listFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := newCmds()

	switch msg := msg.(type) {
	// Keyboard
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, QuitOverallKeys):
			cmds.Add(
				errCmd(&ui.CancelError{}),
			)
		case key.Matches(msg, SelectOverallKeys):
			cmds.Add(
				model.submit(),
			)
		case key.Matches(msg, NextOverallKeys):
			cmds.Add(
				model.updateIndex(model.index.Next()),
			)
		case key.Matches(msg, PreviousOverallKeys):
			cmds.Add(
				model.updateIndex(model.index.Previous()),
			)
		case key.Matches(msg, FirstOverallKeys):
			cmds.Add(
				model.updateIndex(model.index.First()),
			)
		case key.Matches(msg, LastOverallKeys):
			cmds.Add(
				model.updateIndex(model.index.Last()),
			)
		}
	// Mouse
	case tea.MouseMsg:
		for i := range model.items {
			if model.zone.Get(fmt.Sprintf("item_%d", i)).InBounds(msg) {
				if msg.Type == tea.MouseLeft {
					cmds.AddSequence(
						model.updateIndex(i), model.submit(),
					)
				}
				// Hover item
				model.items[i] = cmds.Update(
					model.items[i], hoverMsg(true),
				)
			} else {
				// Not hover item
				model.items[i] = cmds.Update(
					model.items[i], hoverMsg(false),
				)
			}
		}
	}

	// Items
	for i := range model.items {
		model.items[i] = cmds.Update(
			model.items[i],
			focusMsg(i == model.index.Get()),
		)
	}

	return model, cmds.Batch()
}

func (model *listFormModel) updateIndex(index int) tea.Cmd {
	// Set index
	model.setIndex(index)

	// Update form index
	model.form.SetIndex(index)

	// Scroll to item zone
	return toZoneCmd(
		model.zone.Get(fmt.Sprintf("item_%d", index)),
	)
}

func (model *listFormModel) setIndex(index int) {
	model.index.Set(index)
}

func (model *listFormModel) submit() tea.Cmd {
	if err := model.form.Submit(); err != nil {
		return errCmd(err)
	}

	return tea.Quit
}

func (model listFormModel) View() string {
	views := make([]string, len(model.items))
	for i, item := range model.items {
		views[i] = model.zone.Mark(
			fmt.Sprintf("item_%d", i),
			item.View(),
		)
	}

	return model.zone.Scan(
		lipgloss.JoinVertical(lipgloss.Left,
			views...,
		),
	)
}
