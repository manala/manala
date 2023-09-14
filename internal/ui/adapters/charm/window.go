package charm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"math"
	"strings"
)

func newWindowModel(header string, model tea.Model, renderer *lipgloss.Renderer) windowModel {
	return windowModel{
		header: newHeaderModel(header, renderer),
		scroll: scrollModel{
			model: model,
			style: scrollStyle.New(renderer),
		},
	}
}

type windowModel struct {
	width      int
	height     int
	header     tea.Model
	headerView string
	scroll     tea.Model
	err        error
}

func (window windowModel) Init() tea.Cmd {
	return newCmds().
		Init(window.header).
		Init(window.scroll).
		Batch()
}

func (window windowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := newCmds()

	switch msg := msg.(type) {
	// Error
	case error:
		window.err = msg
		return window, tea.Quit
	// Size
	case tea.WindowSizeMsg:
		window.width, window.height = msg.Width, msg.Height
	}

	// Header size
	window.header = cmds.Update(
		window.header, sizeMsg{Width: window.width},
	)

	window.headerView = window.header.View()
	headerHeight := lipgloss.Height(window.headerView)

	// Scroll size
	window.scroll = cmds.Update(
		window.scroll, sizeMsg{
			Width:  window.width,
			Height: max(0, window.height-headerHeight),
		},
	)

	switch msg := msg.(type) {
	// Mouse
	case tea.MouseMsg:
		msg.Y -= headerHeight
		window.scroll = cmds.Update(
			window.scroll, msg,
		)
	default:
		window.scroll = cmds.Update(
			window.scroll, msg,
		)
	}

	return window, cmds.Batch()
}

func (window windowModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		window.headerView,
		window.scroll.View(),
	)
}

func errCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}

type sizeMsg struct {
	Width  int
	Height int
}

/**********/
/* Scroll */
/**********/

type scrollModel struct {
	model  tea.Model
	width  int
	height int
	lines  []string
	line   int
	style  *style
}

func (scroll scrollModel) Init() tea.Cmd {
	return scroll.model.Init()
}

func (scroll scrollModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := newCmds()

	leftFrameSize := scroll.style.GetLeftFrameSize()
	topFrameSize := scroll.style.GetTopFrameSize()

	switch msg := msg.(type) {
	// Size
	case sizeMsg:
		scroll.width = max(2, msg.Width-leftFrameSize-scroll.style.GetRightFrameSize())
		scroll.height = max(1, msg.Height-topFrameSize-scroll.style.GetBottomFrameSize())
		scroll.model = cmds.Update(
			scroll.model, sizeMsg{Width: scroll.width - 1},
		)

	// Mouse
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			scroll.upTo(scroll.line - 1)
		case tea.MouseWheelDown:
			scroll.upTo(scroll.line + 1)
		}
		msg.X -= leftFrameSize
		msg.Y -= topFrameSize - scroll.line
		scroll.model = cmds.Update(
			scroll.model, msg,
		)
	// Zone
	case toZoneMsg:
		scroll.to(msg.StartY, msg.EndY)
	// Everything else
	default:
		scroll.model = cmds.Update(
			scroll.model, msg,
		)
	}

	// Scroll
	scroll.lines = strings.Split(
		scroll.model.View(),
		"\n",
	)

	scroll.upTo(scroll.line)

	return scroll, cmds.Batch()
}

func (scroll *scrollModel) to(lineTop int, lineBottom int) {
	switch {
	case lineTop < scroll.line:
		scroll.upTo(lineTop)
	case lineBottom >= scroll.line+scroll.height:
		scroll.downTo(lineBottom)
	}
}

func (scroll *scrollModel) upTo(line int) {
	if line < 0 {
		scroll.line = 0
	} else if len(scroll.lines) <= scroll.height {
		// Content height inferior to height
		scroll.line = 0
	} else {
		// Content height superior to height
		scroll.line = min(line, len(scroll.lines)-scroll.height)
	}
}

func (scroll *scrollModel) downTo(line int) {
	scroll.line = (line + 1) - scroll.height
}

func (scroll scrollModel) View() string {
	lowBound := scroll.line
	highBound := min(
		len(scroll.lines),
		lowBound+scroll.height,
	)

	view := strings.Join(
		scroll.lines[lowBound:highBound],
		"\n",
	)

	// Bar
	var barView string
	if len(scroll.lines) <= scroll.height {
		barView = strings.TrimRight(
			strings.Repeat(" \n", scroll.height), "\n",
		)
	} else {
		ratio := float64(scroll.height) / float64(len(scroll.lines))
		thumbHeight := max(1, int(math.Round(float64(scroll.height)*ratio)))
		thumbLine := int(math.Round(float64(scroll.line) * ratio))
		barView = strings.TrimRight(
			strings.Repeat("░\n", thumbLine)+
				strings.Repeat("█\n", thumbHeight)+
				strings.Repeat("░\n", max(0, scroll.height-thumbLine-thumbHeight)),
			"\n",
		)
	}

	return scroll.style.Render(
		lipgloss.JoinHorizontal(lipgloss.Top,
			scroll.style.Fit(view, scroll.width-1, scroll.height),
			barView,
		),
	)
}

func toZoneCmd(zone *zone.ZoneInfo) tea.Cmd {
	return func() tea.Msg {
		if zone == nil {
			return nil
		}
		return toZoneMsg(zone)
	}
}

type toZoneMsg *zone.ZoneInfo
