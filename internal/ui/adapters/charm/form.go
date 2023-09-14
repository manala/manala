package charm

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"manala/internal/ui"
	"manala/internal/ui/components"
)

func (adapter *Adapter) Form(header string, form *components.Form) error {
	renderer := adapter.outRenderer

	zone := zone.New()
	defer zone.Close()

	model := formModel{
		fields: make([]tea.Model, len(form.Fields)+1),
		zone:   zone,
	}

	// Build fields
	for i, field := range form.Fields {
		var err error
		model.fields[i], err = newFormFieldModel(field, renderer, zone)
		if err != nil {
			return err
		}
	}

	// Add submit field
	model.fields[len(form.Fields)] = newFormSubmitModel(form, renderer, zone)

	// Focus index
	model.focusIndex = newModelsIndex(&model.fields)

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

type formModel struct {
	fields     []tea.Model
	focusIndex *modelsIndex
	width      int
	zone       *zone.Manager
}

func (model formModel) Init() tea.Cmd {
	return newCmds().
		Init(model.fields...).
		Add(model.Focus(0)).
		Add(textinput.Blink).
		Batch()
}

func (model formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := newCmds()

	switch msg := msg.(type) {
	// Size
	case sizeMsg:
		model.width = msg.Width
	// Keyboard
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, QuitKeys):
			cmds.Add(
				errCmd(&ui.CancelError{}),
			)
		case key.Matches(msg, NextFieldKeys):
			cmds.Add(
				model.Focus(model.focusIndex.Next()),
			)
		case key.Matches(msg, PreviousFieldKeys):
			cmds.Add(
				model.Focus(model.focusIndex.Previous()),
			)
		}
	// Mouse
	case tea.MouseMsg:
		for i := range model.fields {
			if model.zone.Get(fmt.Sprintf("field_%d", i)).InBounds(msg) {
				if msg.Type == tea.MouseLeft {
					// Focus field on left click
					cmds.Add(model.Focus(i))
				}
				// Hover field
				model.fields[i] = cmds.Update(
					model.fields[i], hoverMsg(true),
				)
			} else {
				// Not hover field
				model.fields[i] = cmds.Update(
					model.fields[i], hoverMsg(false),
				)
			}
		}
	// Field
	case formFieldInputMsg:
		cmds.Add(
			model.Focus(model.focusIndex.Next()),
		)
	case formFieldFocusMsg:
		cmds.Add(
			model.Focus(int(msg)),
		)
	}

	// Fields
	for i := range model.fields {
		model.fields[i] = cmds.Update(
			model.fields[i], sizeMsg{Width: model.width}, msg,
		)
	}

	return model, cmds.Batch()
}

func (model *formModel) Focus(index int) tea.Cmd {
	cmds := newCmds()

	curFocusIndex := model.focusIndex.Get()
	if index != curFocusIndex {
		// Blur old field
		model.fields[curFocusIndex] = cmds.Update(
			model.fields[curFocusIndex], focusMsg(false),
		)
	}

	// Update index
	model.focusIndex.Set(index)

	// Focus field
	model.fields[index] = cmds.Update(
		model.fields[index], focusMsg(true),
	)

	// Scroll to field zone
	cmds.Add(
		toZoneCmd(
			model.zone.Get(fmt.Sprintf("field_%d", index)),
		),
	)

	return cmds.Batch()
}

func (model formModel) View() string {
	views := make([]string, len(model.fields))
	for i, field := range model.fields {
		views[i] = model.zone.Mark(
			fmt.Sprintf("field_%d", i),
			field.View(),
		)
	}

	return model.zone.Scan(
		lipgloss.JoinVertical(lipgloss.Left,
			views...,
		),
	)
}
