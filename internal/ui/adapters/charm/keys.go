package charm

import "github.com/charmbracelet/bubbles/key"

var (
	QuitKeys = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
	QuitOverallKeys = key.NewBinding(
		key.WithKeys("ctrl+c", "esc", "q"),
		key.WithHelp("ctrl+c/esc/q", "quit"),
	)
	CancelKeys = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	SelectKeys = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	)
	SelectOverallKeys = key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select"),
	)
	OpenKeys = key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "open"),
	)
	SelectNextKeys = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "select next"),
	)
	SelectPreviousKeys = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "select previous"),
	)
	NextKeys = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next"),
	)
	PreviousKeys = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "previous"),
	)
	NextOverallKeys = key.NewBinding(
		key.WithKeys("tab", "down", "right", "j", "l"),
		key.WithHelp("tab/↓/→/j/l", "next"),
	)
	PreviousOverallKeys = key.NewBinding(
		key.WithKeys("shift+tab", "up", "left", "k", "h"),
		key.WithHelp("shift+tab/↑/←/k/h", "previous"),
	)
	NextFieldKeys = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	)
	PreviousFieldKeys = key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous field"),
	)
	FirstOverallKeys = key.NewBinding(
		key.WithKeys("home", "ctrl+a"),
		key.WithHelp("home/ctrl+a", "first field"),
	)
	LastOverallKeys = key.NewBinding(
		key.WithKeys("end", "ctrl+e"),
		key.WithHelp("end/ctrl+e", "last field"),
	)
)
