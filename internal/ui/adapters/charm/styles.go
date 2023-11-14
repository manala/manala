package charm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor = lipgloss.CompleteAdaptiveColor{
		Light: lipgloss.CompleteColor{ANSI: "8", ANSI256: "233", TrueColor: "#111111"},
		Dark:  lipgloss.CompleteColor{ANSI: "7", ANSI256: "255", TrueColor: "#EEEEEE"},
	}
	primaryDarkColor = lipgloss.CompleteAdaptiveColor{
		Light: lipgloss.CompleteColor{ANSI: "7", ANSI256: "242", TrueColor: "#6B6B6B"},
		Dark:  lipgloss.CompleteColor{ANSI: "8", ANSI256: "246", TrueColor: "#949494"},
	}
	primaryNegativeColor = lipgloss.CompleteAdaptiveColor{
		Light: primaryColor.Dark,
		Dark:  primaryColor.Light,
	}
	secondaryColor = lipgloss.CompleteAdaptiveColor{
		Light: lipgloss.CompleteColor{ANSI: "12", ANSI256: "57", TrueColor: "#5F00FF"},
		Dark:  lipgloss.CompleteColor{ANSI: "6", ANSI256: "79", TrueColor: "#5FD7AF"},
	}
	// Message
	messageColor = primaryDarkColor
	// Levels
	debugColor = primaryColor
	infoColor  = secondaryColor
	warnColor  = lipgloss.CompleteAdaptiveColor{
		Light: lipgloss.CompleteColor{ANSI: "3", ANSI256: "3", TrueColor: "#808000"},
		Dark:  lipgloss.CompleteColor{ANSI: "3", ANSI256: "3", TrueColor: "#808000"},
	}
	errorColor = lipgloss.CompleteAdaptiveColor{
		Light: lipgloss.CompleteColor{ANSI: "1", ANSI256: "1", TrueColor: "#800000"},
		Dark:  lipgloss.CompleteColor{ANSI: "1", ANSI256: "1", TrueColor: "#800000"},
	}
)

var (
	// Levels
	debugStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(debugColor),
	)
	debugSymbolStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(bulletBorder, false, false, false, true).
			BorderForeground(debugColor),
	)
	infoStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(infoColor),
	)
	infoSymbolStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(bulletBorder, false, false, false, true).
			BorderForeground(infoColor),
	)
	warnStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(warnColor),
	)
	warnSymbolStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(warnBorder, false, false, false, true).
			BorderForeground(warnColor),
	)
	errorStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(errorColor),
	)
	errorSymbolStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(crossBorder, false, false, false, true).
			BorderForeground(errorColor),
	)
	// Message
	messageStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1),
	)
	messageMessageStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(messageColor),
	)
	messageAttributesStyle = newStyleDefinition(
		lipgloss.NewStyle(),
	)
	messageAttributeValueStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(messageColor),
	)
	messageDetailsStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingTop(1),
	)
	// Header
	headerStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryColor).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(primaryColor),
	)
	// List
	listItemStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1),
	)
	listItemPrimaryStyle = newStyleDefinition(
		lipgloss.NewStyle().
			MaxHeight(1).
			Foreground(primaryColor).
			PaddingLeft(1).
			Border(bulletBorder, false, false, false, true).
			BorderForeground(primaryColor),
	)
	listItemSecondaryStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryDarkColor).
			PaddingLeft(2),
	)
	// List Form
	listFormItemStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1),
	)
	listFormItemPrimaryStyle = newStyleDefinition(
		lipgloss.NewStyle().
			MaxHeight(1).
			Foreground(primaryColor).
			PaddingLeft(1).
			Border(radioBorder, false, false, false, true).
			BorderForeground(primaryColor),
		withFocusStyle(
			lipgloss.NewStyle().
				MaxHeight(1).
				Foreground(secondaryColor).
				PaddingLeft(1).
				Border(checkedRadioBorder, false, false, false, true).
				BorderForeground(secondaryColor),
		),
		withHoverStyle(
			lipgloss.NewStyle().
				MaxHeight(1).
				Foreground(secondaryColor).
				PaddingLeft(1).
				Border(radioBorder, false, false, false, true).
				BorderForeground(secondaryColor),
		),
	)
	listFormItemSecondaryStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryDarkColor).
			PaddingLeft(2),
	)
	// Animation
	animationStyle = newStyleDefinition(
		lipgloss.NewStyle(),
	)
	// Scroll
	scrollStyle = newStyleDefinition(
		lipgloss.NewStyle(),
	)
	// Form
	formFieldStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(2),
	)
	formLabelStyle = newStyleDefinition(
		lipgloss.NewStyle().
			MaxHeight(1).
			Foreground(primaryColor),
	)
	formHelpStyle = newStyleDefinition(
		lipgloss.NewStyle().
			MaxHeight(1).
			Foreground(primaryDarkColor),
	)
	formViolationStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(errorColor),
	)
	formViolationSymbolStyle = newStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(crossBorder, false, false, false, true).
			BorderForeground(errorColor),
	)
	formTextStyle = newStyleDefinition(
		lipgloss.NewStyle(),
	)
	formTextInputStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryDarkColor).
			Background(primaryNegativeColor),
		withFocusStyle(
			lipgloss.NewStyle().
				Foreground(secondaryColor).
				Background(primaryNegativeColor),
		),
		withHoverStyle(
			lipgloss.NewStyle().
				Foreground(primaryNegativeColor).
				Background(secondaryColor),
		),
	)
	formTextInputCursorStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor).
			Background(primaryNegativeColor),
	)
	formTextInputCursorTextStyle = formTextInputStyle
	formSelectStyle              = newStyleDefinition(
		lipgloss.NewStyle().
			MaxHeight(1).
			Foreground(primaryDarkColor).
			PaddingLeft(1).
			Border(whiteRightTriangleBorder, false, false, false, true).
			BorderForeground(primaryDarkColor),
		withFocusStyle(
			lipgloss.NewStyle().
				MaxHeight(1).
				Foreground(secondaryColor).
				PaddingLeft(1).
				Border(blackRightTriangleBorder, false, false, false, true).
				BorderForeground(secondaryColor),
		),
		withHoverStyle(
			lipgloss.NewStyle().
				MaxHeight(1).
				Foreground(secondaryColor).
				PaddingLeft(1).
				Border(whiteRightTriangleBorder, false, false, false, true).
				BorderForeground(secondaryColor),
		),
	)
	formSelectOptionStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryDarkColor).
			PaddingLeft(2),
		withFocusStyle(
			lipgloss.NewStyle().
				Foreground(secondaryColor).
				PaddingLeft(1).
				Border(blackRightTriangleBorder, false, false, false, true).
				BorderForeground(secondaryColor),
		),
		withFocusHoverStyle(
			lipgloss.NewStyle().
				Foreground(primaryNegativeColor).
				Background(secondaryColor).
				PaddingLeft(1).
				Border(blackRightTriangleBorder, false, false, false, true).
				BorderForeground(primaryNegativeColor).
				BorderBackground(secondaryColor),
		),
		withHoverStyle(
			lipgloss.NewStyle().
				Foreground(primaryNegativeColor).
				Background(secondaryColor).
				PaddingLeft(2),
		),
	)
	formSubmitStyle = newStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryNegativeColor).
			Background(primaryColor).
			Margin(1, 0, 0, 2).
			Padding(0, 1, 0, 1),
		withFocusStyle(
			lipgloss.NewStyle().
				Foreground(primaryNegativeColor).
				Background(secondaryColor).
				Margin(1, 0, 0, 2).
				Padding(0, 1, 0, 1),
		),
		withHoverStyle(
			lipgloss.NewStyle().
				Foreground(primaryNegativeColor).
				Background(secondaryColor).
				Margin(1, 0, 0, 2).
				Padding(0, 1, 0, 1),
		),
	)
)

var (
	bulletBorder = lipgloss.Border{
		Left: "•",
	}
	radioBorder = lipgloss.Border{
		Left: "○",
	}
	checkedRadioBorder = lipgloss.Border{
		Left: "●",
	}
	warnBorder = lipgloss.Border{
		Left: "⚠",
	}
	crossBorder = lipgloss.Border{
		Left: "⨯",
	}
	whiteRightTriangleBorder = lipgloss.Border{
		Left: "▷",
	}
	blackRightTriangleBorder = lipgloss.Border{
		Left: "▶",
	}
)

func newStyleDefinition(style lipgloss.Style, opts ...styleOption) *styleDefinition {
	definition := &styleDefinition{
		style: style,
	}

	// Apply options
	for _, opt := range opts {
		opt(definition)
	}

	return definition
}

type styleOption func(definition *styleDefinition)

func withFocusStyle(style lipgloss.Style) styleOption {
	return func(definition *styleDefinition) {
		definition.focusStyle = &style
	}
}

func withFocusHoverStyle(style lipgloss.Style) styleOption {
	return func(definition *styleDefinition) {
		definition.focusHoverStyle = &style
	}
}

func withHoverStyle(style lipgloss.Style) styleOption {
	return func(definition *styleDefinition) {
		definition.hoverStyle = &style
	}
}

type styleDefinition struct {
	style           lipgloss.Style
	focusStyle      *lipgloss.Style
	focusHoverStyle *lipgloss.Style
	hoverStyle      *lipgloss.Style
}

func (definition *styleDefinition) New(renderer *lipgloss.Renderer) *style {
	style := &style{
		definition: definition,
		renderer:   renderer,
	}
	style.update()
	return style
}

type style struct {
	definition *styleDefinition
	renderer   *lipgloss.Renderer
	style      lipgloss.Style
	focus      bool
	hover      bool
}

func (s *style) Update(msg tea.Msg) {
	var changed bool

	switch _msg := msg.(type) {
	case focusMsg:
		changed = s.focus != bool(_msg)
		s.focus = bool(_msg)
	case hoverMsg:
		changed = s.hover != bool(_msg)
		s.hover = bool(_msg)
	}

	if changed {
		s.update()
	}
}

func (s *style) update() {
	style := s.definition.style

	switch {
	case s.focus && s.hover && s.definition.focusHoverStyle != nil:
		style = *s.definition.focusHoverStyle
	case s.focus && s.definition.focusStyle != nil:
		style = *s.definition.focusStyle
	case s.hover && s.definition.hoverStyle != nil:
		style = *s.definition.hoverStyle
	}

	s.style = style.Copy().Renderer(s.renderer)
}

func (s *style) Style() lipgloss.Style {
	return s.style
}

func (s *style) Render(strs ...string) string {
	return s.style.Render(strs...)
}

func (s *style) GetTopFrameSize() int {
	return s.style.GetMarginTop() +
		s.style.GetPaddingTop() +
		s.style.GetBorderTopSize()
}

func (s *style) GetRightFrameSize() int {
	return s.style.GetMarginRight() +
		s.style.GetPaddingRight() +
		s.style.GetBorderRightSize()
}

func (s *style) GetBottomFrameSize() int {
	return s.style.GetMarginBottom() +
		s.style.GetPaddingBottom() +
		s.style.GetBorderBottomSize()
}

func (s *style) GetLeftFrameSize() int {
	return s.style.GetMarginLeft() +
		s.style.GetPaddingLeft() +
		s.style.GetBorderLeftSize()
}

func (s *style) GetHorizontalFrameSize() int {
	return s.style.GetHorizontalFrameSize()
}

func (s *style) GetVerticalFrameSize() int {
	return s.style.GetVerticalFrameSize()
}

func (s *style) Fit(str string, width int, height int) string {
	style := s.style.Copy()

	if width > 0 {
		if lipgloss.Width(str) < width {
			style.Width(width)
		} else {
			style = style.MaxWidth(width)
		}
	}

	if height > 0 {
		if lipgloss.Height(str) < height {
			style.Height(height)
		} else {
			style = style.MaxHeight(height)
		}
	}

	return style.
		UnsetMargins().
		UnsetPadding().
		UnsetBorderStyle().
		Render(str)
}
