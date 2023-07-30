package lipgloss

import (
	"github.com/charmbracelet/lipgloss"
	"manala/internal/ui/components"
)

func (output *Output) Table(table *components.Table) {
	// Get primary width
	primaryWidth := 0
	for _, row := range table.Rows {
		width := lipgloss.Width(row.Primary)
		if width > primaryWidth {
			primaryWidth = width
		}
	}

	// Styles
	primaryStyle := output.outStyle().
		Foreground(color).
		Bold(true).
		Width(primaryWidth).
		MarginRight(2)
	secondaryStyle := output.outStyle().
		Foreground(color).
		Italic(true)

	for _, row := range table.Rows {
		output.writeOutString(
			primaryStyle.Render(row.Primary) +
				secondaryStyle.Render(row.Secondary) +
				"\n",
		)
	}
}
