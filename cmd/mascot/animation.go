package mascot

import (
	_ "embed"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand/v2"
	"time"
)

type animation struct {
	duration  int
	repeat    int
	style     lipgloss.Style
	frame     *string
	frameYell *string
	yell      bool
	width     int
	height    int
	err       error
}

func (m animation) Init() tea.Cmd {
	return tea.Sequence(
		tea.SetWindowTitle("Quack Quack"),
		func() tea.Msg {
			if m.repeat == 0 {
				return animationStopMsg{}
			}

			return animationYellMsg(true)
		},
	)
}

func (m animation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var c tea.Cmd

	switch msg := msg.(type) {
	// Size
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	// Error
	case error:
		m.err = msg
		c = tea.Quit
	// Keys
	case tea.KeyMsg:
		switch msg.String() {
		// Keys - Quit
		case "ctrl+c", "esc", "q":
			c = tea.Quit
		}
	// Animation - Stop
	case animationStopMsg:
		c = tea.Quit
	// Animation - Yell
	case animationYellMsg:
		m.yell = bool(msg)
		if msg {
			c = m.yellStart
		} else {
			if m.repeat > 0 {
				m.repeat--
			}
			c = m.yellStop
		}
	}

	return m, c
}

func (m animation) View() string {
	frame := m.frame
	if m.yell {
		frame = m.frameYell
	}

	// Render
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.style.
			MaxWidth(m.width).
			MaxHeight(m.height).
			Render(*frame),
	)
}

func (m animation) yellStart() tea.Msg {
	// Has yell audio ?
	var modelInterface any = m
	if _model, ok := modelInterface.(interface{ yellAudio() error }); ok {
		if err := _model.yellAudio(); err != nil {
			return err
		}
	} else {
		m.pause()
	}

	// Yell is finished, stop it
	return animationYellMsg(false)
}

func (m animation) yellStop() tea.Msg {
	m.pause()

	if m.repeat == 0 {
		return animationStopMsg{}
	}

	return animationYellMsg(true)
}

func (m animation) pause() {
	duration := (m.duration / 2) + rand.IntN(m.duration)
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

type animationStopMsg struct{}
type animationYellMsg bool
