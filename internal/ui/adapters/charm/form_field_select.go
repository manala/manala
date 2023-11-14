package charm

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"manala/internal/ui/components"
)

func newFormFieldSelectModel(field *components.FormFieldSelect, renderer *lipgloss.Renderer, zone *zone.Manager) formFieldSelectModel {
	model := formFieldSelectModel{
		field:       field,
		open:        false,
		options:     make([]tea.Model, len(field.Options)),
		zone:        zone,
		zonePrefix:  zone.NewPrefix(),
		style:       formFieldStyle.New(renderer),
		labelStyle:  formLabelStyle.New(renderer),
		helpStyle:   formHelpStyle.New(renderer),
		selectStyle: formSelectStyle.New(renderer),
		formFieldModel: &formFieldModel{
			violationStyle:       formViolationStyle.New(renderer),
			violationSymbolStyle: formViolationSymbolStyle.New(renderer),
		},
	}

	// Build options
	var labelWidth int
	for i, option := range field.Options {
		labelWidth = max(labelWidth, lipgloss.Width(option.Label()))
		model.options[i] = newFormFieldSelectOptionModel(
			option,
			&labelWidth,
			renderer,
		)
	}

	// Indexes
	model.index = newModelsIndex(&model.options)
	model.hoverIndex = newModelsIndex(&model.options)
	model.hoverIndex.Circular = true

	return model
}

type formFieldSelectModel struct {
	field       *components.FormFieldSelect
	focus       bool
	open        bool
	options     []tea.Model
	index       *modelsIndex
	hoverIndex  *modelsIndex
	zone        *zone.Manager
	zonePrefix  string
	style       *style
	labelStyle  *style
	helpStyle   *style
	selectStyle *style
	*formFieldModel
}

func (model formFieldSelectModel) Init() tea.Cmd {
	cmds := newCmds().
		Init(model.options...)

	// Initial index
	index := model.field.GetIndex()
	if index == -1 {
		cmds.Add(
			model.updateIndex(0),
		)
	} else {
		model.setIndex(index)
	}

	return cmds.Batch()
}

func (model formFieldSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := newCmds()

	// Update styles
	model.style.Update(msg)
	model.labelStyle.Update(msg)
	model.helpStyle.Update(msg)
	model.selectStyle.Update(msg)

	switch msg := msg.(type) {
	// Focus
	case focusMsg:
		model.focus = bool(msg)
		// Force closing on blur
		if !msg {
			model.Close()
		}
	// Keys
	case tea.KeyMsg:
		// Bypass keys events when field don't have focus
		if !model.focus {
			return model, nil
		}
		if model.open {
			switch {
			case key.Matches(msg, NextKeys):
				model.hoverIndex.SetNext()
			case key.Matches(msg, PreviousKeys):
				model.hoverIndex.SetPrevious()
			case key.Matches(msg, CancelKeys):
				model.Close()
			case key.Matches(msg, SelectKeys):
				cmds.Add(
					model.updateIndex(model.hoverIndex.Get()), formFieldInput,
				)
			}
		} else {
			switch {
			case key.Matches(msg, SelectNextKeys):
				cmds.Add(
					model.updateIndex(model.index.Next()),
				)
			case key.Matches(msg, SelectPreviousKeys):
				cmds.Add(
					model.updateIndex(model.index.Previous()),
				)
			case key.Matches(msg, OpenKeys),
				key.Matches(msg, NextKeys),
				key.Matches(msg, PreviousKeys):
				model.Open(false)
			case key.Matches(msg, SelectKeys):
				cmds.Add(formFieldInput)
			}
		}
	// Mouse
	case tea.MouseMsg:
		if model.open {
			// Options
			for i := range model.options {
				if model.zone.Get(fmt.Sprintf("%soption_%d", model.zonePrefix, i)).InBounds(msg) {
					model.hoverIndex.Set(i)
					if msg.Type == tea.MouseLeft {
						cmds.Add(model.updateIndex(i))
						model.Close()
					}
				}
			}
		} else {
			switch msg.Type {
			case tea.MouseLeft:
				// Open on left click
				if model.zone.Get(model.zonePrefix + "select").InBounds(msg) {
					model.Open(true)
				}
			}
		}
	}

	// Options
	for i := range model.options {
		model.options[i] = cmds.Update(
			model.options[i],
			focusMsg(i == model.index.Get()),
			hoverMsg(i == model.hoverIndex.Get()),
		)
	}

	return model, cmds.Batch()
}

func (model *formFieldSelectModel) updateIndex(index int) tea.Cmd {
	// Set index
	model.setIndex(index)

	// Update field index
	if err := model.field.SetIndex(index); err != nil {
		return errCmd(err)
	}

	return nil
}

func (model *formFieldSelectModel) setIndex(index int) {
	model.index.Set(index)
	model.hoverIndex.Set(index)
}

func (model *formFieldSelectModel) Open(resetHover bool) {
	if resetHover {
		model.hoverIndex.Reset()
	}
	model.open = true
}

func (model *formFieldSelectModel) Close() {
	model.open = false
}

func (model formFieldSelectModel) View() string {
	// Label
	view := model.labelStyle.Render(
		model.field.Label(),
	)

	// Help
	help := model.field.Help()
	if help != "" {
		view = lipgloss.JoinVertical(lipgloss.Left,
			view,
			model.helpStyle.Render(help),
		)
	}

	if model.open {
		views := make([]string, len(model.options))
		for i, option := range model.options {
			views[i] = model.zone.Mark(fmt.Sprintf("%soption_%d", model.zonePrefix, i),
				option.View(),
			)
		}

		view = lipgloss.JoinVertical(lipgloss.Left,
			view,
			lipgloss.JoinVertical(lipgloss.Left, views...),
		)
	} else {
		view = lipgloss.JoinVertical(lipgloss.Left,
			view,
			model.zone.Mark(model.zonePrefix+"select",
				model.selectStyle.Render(
					model.field.Options[model.index.Get()].Label(),
				),
			),
		)
	}

	// Violations
	if violationsView := model.formFieldModel.viewViolations(model.field.Violations()); violationsView != "" {
		view = lipgloss.JoinVertical(lipgloss.Left, view, violationsView)
	}

	return model.style.Render(view)
}

/**********/
/* Option */
/**********/

func newFormFieldSelectOptionModel(option *components.FormFieldSelectOption, width *int, renderer *lipgloss.Renderer) formFieldSelectOptionModel {
	return formFieldSelectOptionModel{
		option: option,
		width:  width,
		style:  formSelectOptionStyle.New(renderer),
	}
}

type formFieldSelectOptionModel struct {
	option *components.FormFieldSelectOption
	width  *int
	style  *style
}

func (model formFieldSelectOptionModel) Init() tea.Cmd {
	return nil
}

func (model formFieldSelectOptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update styles
	model.style.Update(msg)

	return model, nil
}

func (model formFieldSelectOptionModel) View() string {
	return model.style.Render(
		model.style.Fit(
			model.option.Label(), *model.width, 0,
		),
	)
}
