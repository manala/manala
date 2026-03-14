package cmd

import (
	"charm.land/lipgloss/v2"
)

var Styles = struct {
	Primary   lipgloss.Style
	Secondary lipgloss.Style
}{
	Primary:   lipgloss.NewStyle().Foreground(lipgloss.NoColor{}),
	Secondary: lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(247)),
}
