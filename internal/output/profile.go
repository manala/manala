package output

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/gdamore/tcell/v3/color"
)

var Plain = Profile{profile: colorprofile.ASCII}

type Profile struct {
	light   bool
	profile colorprofile.Profile
}

func (p Profile) Rich() bool     { return p.profile > colorprofile.ASCII }
func (p Profile) Basic() bool    { return p.profile == colorprofile.ANSI }
func (p Profile) Extended() bool { return p.profile == colorprofile.ANSI256 }
func (p Profile) True() bool     { return p.profile == colorprofile.TrueColor }
func (p Profile) Light() bool    { return p.light }

// Style represents the default style.
func (p Profile) Style() lipgloss.Style {
	style := lipgloss.NewStyle()
	if p.Rich() {
		style = style.Foreground(lipgloss.NoColor{})
	}
	return style
}

// Color represents the default color.
func (p Profile) Color() color.Color {
	return color.Default
}

func (p Profile) ReverseColor() color.Color {
	return color.Black
}

func (p Profile) MutedStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if p.Rich() {
		style = style.Faint(true)
	}
	return style
}

func (p Profile) MutedColor() color.Color {
	// Mimic faint
	return color.XTerm247
}

func (p Profile) EmphasisColor() color.Color {
	return color.White
}

func (p Profile) LitteralStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if p.Rich() {
		style = style.Foreground(lipgloss.Cyan)
	}
	return style
}

func (p Profile) ErrorStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if p.Rich() {
		style = style.Foreground(lipgloss.Red)
	}
	return style
}

func (p Profile) ErrorColor() color.Color {
	return color.Maroon
}

func (p Profile) WarnStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if p.Rich() {
		style = style.Foreground(lipgloss.Yellow)
	}
	return style
}

func (p Profile) InfoStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if p.Rich() {
		style = style.Foreground(lipgloss.Green)
	}
	return style
}

func (p Profile) DebugStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if p.Rich() {
		style = style.Foreground(lipgloss.Blue)
	}
	return style
}
