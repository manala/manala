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
	// Message.
	messageColor = primaryDarkColor
	// Levels.
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
	// Levels.
	debugStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(debugColor),
	)
	debugSymbolStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(bulletBorder, false, false, false, true).
			BorderForeground(debugColor),
	)
	infoStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(infoColor),
	)
	infoSymbolStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(bulletBorder, false, false, false, true).
			BorderForeground(infoColor),
	)
	warnStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(warnColor),
	)
	warnSymbolStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(warnBorder, false, false, false, true).
			BorderForeground(warnColor),
	)
	errorStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(errorColor),
	)
	errorSymbolStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(crossBorder, false, false, false, true).
			BorderForeground(errorColor),
	)
	// Message.
	messageStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1),
	)
	messageMessageStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Width(32).
			Foreground(messageColor),
	)
	messageAttributesStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1),
	)
	messageAttributeValueStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(messageColor),
	)
	messageDetailsStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingTop(1),
	)
	// Header.
	headerStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryColor).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(primaryColor),
	)
	// List.
	listItemStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1),
	)
	listItemPrimaryStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			MaxHeight(1).
			Foreground(primaryColor).
			PaddingLeft(1).
			Border(bulletBorder, false, false, false, true).
			BorderForeground(primaryColor),
	)
	listItemSecondaryStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryDarkColor).
			PaddingLeft(2),
	)
	// List Form.
	listFormItemStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1),
	)
	listFormItemPrimaryStyle = NewStyleDefinition(
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
	listFormItemSecondaryStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(primaryDarkColor).
			PaddingLeft(2),
	)
	// Scroll.
	scrollStyle = NewStyleDefinition(
		lipgloss.NewStyle(),
	)
	// Form.
	formFieldStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(2),
	)
	formLabelStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			MaxHeight(1).
			Foreground(primaryColor),
	)
	formHelpStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			MaxHeight(1).
			Foreground(primaryDarkColor),
	)
	formViolationStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Foreground(errorColor),
	)
	formViolationSymbolStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			PaddingLeft(1).
			Border(crossBorder, false, false, false, true).
			BorderForeground(errorColor),
	)
	formTextStyle = NewStyleDefinition(
		lipgloss.NewStyle(),
	)
	formTextInputStyle = NewStyleDefinition(
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
	formTextInputCursorStyle = NewStyleDefinition(
		lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor).
			Background(primaryNegativeColor),
	)
	formTextInputCursorTextStyle = formTextInputStyle
	formSelectStyle              = NewStyleDefinition(
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
	formSelectOptionStyle = NewStyleDefinition(
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
	formSubmitStyle = NewStyleDefinition(
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

type styleOption func(definition *StyleDefinition)

func withFocusStyle(style lipgloss.Style) styleOption {
	return func(definition *StyleDefinition) {
		definition.focusStyle = &style
	}
}

func withFocusHoverStyle(style lipgloss.Style) styleOption {
	return func(definition *StyleDefinition) {
		definition.focusHoverStyle = &style
	}
}

func withHoverStyle(style lipgloss.Style) styleOption {
	return func(definition *StyleDefinition) {
		definition.hoverStyle = &style
	}
}

type StyleDefinition struct {
	style           lipgloss.Style
	focusStyle      *lipgloss.Style
	focusHoverStyle *lipgloss.Style
	hoverStyle      *lipgloss.Style
}

func NewStyleDefinition(style lipgloss.Style, opts ...styleOption) *StyleDefinition {
	definition := &StyleDefinition{
		style: style,
	}

	// Apply options
	for _, opt := range opts {
		opt(definition)
	}

	return definition
}

func (definition *StyleDefinition) New(renderer *lipgloss.Renderer) *Style {
	style := &Style{
		definition: definition,
		renderer:   renderer,
	}
	style.update()

	return style
}

type Style struct {
	definition *StyleDefinition
	renderer   *lipgloss.Renderer
	style      lipgloss.Style
	focus      bool
	hover      bool
}

func (s *Style) Update(msg tea.Msg) {
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

func (s *Style) Style() lipgloss.Style {
	return s.style
}

func (s *Style) Render(strs ...string) string {
	return s.style.Render(strs...)
}

func (s *Style) GetTopFrameSize() int {
	return s.style.GetMarginTop() +
		s.style.GetPaddingTop() +
		s.style.GetBorderTopSize()
}

func (s *Style) GetRightFrameSize() int {
	return s.style.GetMarginRight() +
		s.style.GetPaddingRight() +
		s.style.GetBorderRightSize()
}

func (s *Style) GetBottomFrameSize() int {
	return s.style.GetMarginBottom() +
		s.style.GetPaddingBottom() +
		s.style.GetBorderBottomSize()
}

func (s *Style) GetLeftFrameSize() int {
	return s.style.GetMarginLeft() +
		s.style.GetPaddingLeft() +
		s.style.GetBorderLeftSize()
}

func (s *Style) GetHorizontalFrameSize() int {
	return s.style.GetHorizontalFrameSize()
}

func (s *Style) GetVerticalFrameSize() int {
	return s.style.GetVerticalFrameSize()
}

func (s *Style) Fit(str string, width, height int) string {
	style := s.style

	if width > 0 {
		if lipgloss.Width(str) < width {
			style = style.Width(width)
		} else {
			style = style.MaxWidth(width)
		}
	}

	if height > 0 {
		if lipgloss.Height(str) < height {
			style = style.Height(height)
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

func (s *Style) update() {
	style := s.definition.style

	switch {
	case s.focus && s.hover && s.definition.focusHoverStyle != nil:
		style = *s.definition.focusHoverStyle
	case s.focus && s.definition.focusStyle != nil:
		style = *s.definition.focusStyle
	case s.hover && s.definition.hoverStyle != nil:
		style = *s.definition.hoverStyle
	}

	s.style = style.Renderer(s.renderer)
}
