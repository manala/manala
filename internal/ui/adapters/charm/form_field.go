package charm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"manala/internal/serrors"
	"manala/internal/ui/components"
	"manala/internal/validator"
)

func newFormFieldModel(field components.FormField, renderer *lipgloss.Renderer, zone *zone.Manager) (tea.Model, error) {
	switch field := field.(type) {
	case *components.FormFieldText:
		return newFormFieldTextModel(field, renderer, zone), nil
	case *components.FormFieldSelect:
		return newFormFieldSelectModel(field, renderer, zone), nil
	}

	return nil, serrors.New("unknown form field type")
}

type formFieldModel struct {
	violationStyle       *style
	violationSymbolStyle *style
}

func (model formFieldModel) viewViolations(violations validator.Violations) string {
	views := make([]string, len(violations))
	for i := range violations {
		views[i] = lipgloss.JoinHorizontal(lipgloss.Top,
			model.violationSymbolStyle.Render(""),
			model.violationStyle.Render(violations[i].Message),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func formFieldInput() tea.Msg {
	return formFieldInputMsg{}
}

// Field has just got user input
type formFieldInputMsg struct{}

func formFieldFocusCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return formFieldFocusMsg(index)
	}
}

// Focus on field index
type formFieldFocusMsg int
