package charm

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"manala/internal/ui/components"
)

func newFormSubmitModel(form *components.Form, renderer *lipgloss.Renderer, zone *zone.Manager) formSubmitModel {
	return formSubmitModel{
		form:       form,
		zone:       zone,
		zonePrefix: zone.NewPrefix(),
		style:      formSubmitStyle.New(renderer),
	}
}

type formSubmitModel struct {
	form       *components.Form
	focus      bool
	zone       *zone.Manager
	zonePrefix string
	style      *style
}

func (model formSubmitModel) Init() tea.Cmd {
	return nil
}

func (model formSubmitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := newCmds()

	// Update styles
	model.style.Update(msg)

	switch msg := msg.(type) {
	// Focus
	case focusMsg:
		model.focus = bool(msg)
	// Keys
	case tea.KeyMsg:
		// Bypass keys events when field don't have focus
		if !model.focus {
			return model, nil
		}
		switch {
		case key.Matches(msg, SelectKeys):
			cmds.Add(model.submit())
		}
	// Mouse
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			// Submit on left click
			if model.zone.Get(model.zonePrefix + "submit").InBounds(msg) {
				cmds.Add(model.submit())
			}
		}
	}

	return model, cmds.Batch()
}

func (model *formSubmitModel) submit() tea.Cmd {
	ok, err := model.form.Submit()
	if err != nil {
		return errCmd(err)
	}
	if !ok {
		// Find thirst field with violations
		var index int
		for i, field := range model.form.Fields {
			if len(field.Violations()) != 0 {
				index = i
				break
			}
		}
		return formFieldFocusCmd(index)
	}

	return tea.Quit
}

func (model formSubmitModel) View() string {
	return model.zone.Mark(model.zonePrefix+"submit",
		model.style.Render("Submit"),
	)
}
