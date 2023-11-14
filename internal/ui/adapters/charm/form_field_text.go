package charm

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"manala/internal/ui/components"
)

func newFormFieldTextModel(field *components.FormFieldText, renderer *lipgloss.Renderer, zone *zone.Manager) formFieldTextModel {
	input := textinput.New()
	input.Prompt = ""

	// Max length
	input.CharLimit = field.MaxLength

	model := formFieldTextModel{
		field:                    field,
		value:                    new(string),
		input:                    input,
		zone:                     zone,
		zonePrefix:               zone.NewPrefix(),
		style:                    formFieldStyle.New(renderer),
		labelStyle:               formLabelStyle.New(renderer),
		helpStyle:                formHelpStyle.New(renderer),
		textStyle:                formTextStyle.New(renderer),
		textInputStyle:           formTextInputStyle.New(renderer),
		textInputCursorStyle:     formTextInputCursorStyle.New(renderer),
		textInputCursorTextStyle: formTextInputCursorTextStyle.New(renderer),
		formFieldModel: &formFieldModel{
			violationStyle:       formViolationStyle.New(renderer),
			violationSymbolStyle: formViolationSymbolStyle.New(renderer),
		},
	}

	// Initial value
	if value, ok := field.Get().(string); ok {
		model.input.SetValue(value)
		model.setValue(value)
	}

	return model
}

type formFieldTextModel struct {
	field                    *components.FormFieldText
	value                    *string
	input                    textinput.Model
	focus                    bool
	width                    int
	zone                     *zone.Manager
	zonePrefix               string
	style                    *style
	labelStyle               *style
	helpStyle                *style
	textStyle                *style
	textInputStyle           *style
	textInputCursorStyle     *style
	textInputCursorTextStyle *style
	*formFieldModel
}

func (model formFieldTextModel) Init() tea.Cmd {
	return nil
}

func (model formFieldTextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := newCmds()

	// Update styles
	model.style.Update(msg)
	model.labelStyle.Update(msg)
	model.helpStyle.Update(msg)
	model.textStyle.Update(msg)
	model.textInputStyle.Update(msg)
	model.textInputCursorStyle.Update(msg)
	model.textInputCursorTextStyle.Update(msg)

	switch msg := msg.(type) {
	// Size
	case sizeMsg:
		model.width = msg.Width
	// Focus
	case focusMsg:
		model.focus = bool(msg)
		if msg {
			cmds.Add(model.input.Focus())
		} else {
			model.input.Blur()
		}
	// Keys
	case tea.KeyMsg:
		// Bypass keys events when field don't have focus
		if !model.focus {
			return model, nil
		}
		switch {
		case key.Matches(msg, SelectKeys):
			cmds.Add(formFieldInput)
		}
	// Mouse
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			zone := model.zone.Get(model.zonePrefix + "input")
			// Set cursor position on left click
			if zone.InBounds(msg) {
				x, _ := zone.Pos(msg)
				model.input.SetCursor(x)
			}
		}
	}

	// Input styles
	model.input.TextStyle = model.textInputStyle.Style()
	model.input.Cursor.Style = model.textInputCursorStyle.Style()
	model.input.Cursor.TextStyle = model.textInputCursorTextStyle.Style()

	width := model.width -
		model.style.GetHorizontalFrameSize() -
		model.textStyle.GetHorizontalFrameSize()

	if model.input.CharLimit != 0 {
		model.input.Width = min(model.input.CharLimit, width)
	} else {
		model.input.Width = min(30, width)
	}

	// Input
	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)
	cmds.Add(cmd)

	// Field
	cmds.Add(
		model.updateValue(model.input.Value()),
	)

	return model, cmds.Batch()
}

func (model *formFieldTextModel) updateValue(value string) tea.Cmd {
	if *model.value == value {
		return nil
	}

	// Set value
	model.setValue(value)

	// Update field
	if err := model.field.Set(value); err != nil {
		return errCmd(err)
	}

	return nil
}

func (model *formFieldTextModel) setValue(value string) {
	*model.value = value
}

func (model formFieldTextModel) View() string {
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

	// Text
	view = lipgloss.JoinVertical(lipgloss.Left,
		view,
		model.textStyle.Render(
			model.zone.Mark(model.zonePrefix+"input",
				model.input.View(),
			),
		),
	)

	// Violations
	if violationsView := model.formFieldModel.viewViolations(model.field.Violations()); violationsView != "" {
		view = lipgloss.JoinVertical(lipgloss.Left, view, violationsView)
	}

	return model.style.Render(view)
}
