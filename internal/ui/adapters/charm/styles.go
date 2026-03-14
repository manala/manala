package charm

import (
	"charm.land/lipgloss/v2"
)

var (
	primaryColor     = lipgloss.ANSIColor(255)
	primaryDarkColor = lipgloss.ANSIColor(246)
	secondaryColor   = lipgloss.ANSIColor(79)
	// Message.
	messageColor = primaryDarkColor
	// Levels.
	debugColor = primaryColor
	infoColor  = secondaryColor
	warnColor  = lipgloss.ANSIColor(3)
	errorColor = lipgloss.ANSIColor(1)
)

var (
	// Levels.
	debugStyle = lipgloss.NewStyle().
			Foreground(debugColor)
	debugSymbolStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Border(bulletBorder, false, false, false, true).
				BorderForeground(debugColor)
	infoStyle = lipgloss.NewStyle().
			Foreground(infoColor)
	infoSymbolStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Border(bulletBorder, false, false, false, true).
			BorderForeground(infoColor)
	warnStyle = lipgloss.NewStyle().
			Foreground(warnColor)
	warnSymbolStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Border(warnBorder, false, false, false, true).
			BorderForeground(warnColor)
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor)
	errorSymbolStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Border(crossBorder, false, false, false, true).
				BorderForeground(errorColor)
	// Message.
	messageStyle = lipgloss.NewStyle().
			PaddingLeft(1)
	messageMessageStyle = lipgloss.NewStyle().
				Width(32).
				Foreground(messageColor)
	messageAttributesStyle = lipgloss.NewStyle().
				PaddingLeft(1)
	messageAttributeValueStyle = lipgloss.NewStyle().
					Foreground(messageColor)
	messageDetailsStyle = lipgloss.NewStyle().
				PaddingTop(1)
)

var (
	bulletBorder = lipgloss.Border{
		Left: "•",
	}
	warnBorder = lipgloss.Border{
		Left: "⚠",
	}
	crossBorder = lipgloss.Border{
		Left: "⨯",
	}
)
