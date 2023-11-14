package charm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newHeaderModel(header string, renderer *lipgloss.Renderer) headerModel {
	return headerModel{
		header: header,
		style:  headerStyle.New(renderer),
	}

}

type headerModel struct {
	header string
	width  int
	style  *style
}

func (header headerModel) Init() tea.Cmd {
	return nil
}

func (header headerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Size
	case sizeMsg:
		header.width = max(1, msg.Width-header.style.GetHorizontalFrameSize())
	}

	return header, nil
}

func (header headerModel) View() string {
	return header.style.Render(
		header.style.Fit(header.header, header.width, 0),
	)
}
