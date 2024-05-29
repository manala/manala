package mascot

import (
	_ "embed"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand/v2"
	"time"
)

type Animation struct {
	Title     string
	Duration  int
	Repeat    int
	Style     lipgloss.Style
	QuitKey   key.Binding
	Frame     *string
	FrameYell *string
	AudioYell *[]byte
	Err       error
	width     int
	height    int
	yell      bool
}

func (model Animation) Init() tea.Cmd {
	return tea.Sequence(
		tea.SetWindowTitle(model.Title),
		func() tea.Msg {
			if model.Repeat == 0 {
				return animationStopMsg{}
			}

			return animationYellMsg(true)
		},
	)
}

func (model Animation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Size
	case tea.WindowSizeMsg:
		model.width, model.height = msg.Width, msg.Height
	// Error
	case error:
		model.Err = msg
		cmd = tea.Quit
	// Keys
	case tea.KeyMsg:
		// Quit
		if key.Matches(msg, model.QuitKey) {
			cmd = tea.Quit
		}
	// Animation - Stop
	case animationStopMsg:
		cmd = tea.Quit
	// Animation - Yell
	case animationYellMsg:
		model.yell = bool(msg)
		if msg {
			cmd = model.yellStart
		} else {
			if model.Repeat > 0 {
				model.Repeat--
			}
			cmd = model.yellStop
		}
	}

	return model, cmd
}

func (model Animation) View() string {
	frame := model.Frame
	if model.yell {
		frame = model.FrameYell
	}

	// Render
	return lipgloss.Place(
		model.width, model.height,
		lipgloss.Center, lipgloss.Center,
		model.Style.
			MaxWidth(model.width).
			MaxHeight(model.height).
			Render(*frame),
	)
}

func (model Animation) yellStart() tea.Msg {
	// Has yell audio ?
	var modelInterface any = model
	if _model, ok := modelInterface.(interface{ yellAudio() error }); ok {
		if err := _model.yellAudio(); err != nil {
			return err
		}
	} else {
		model.pause()
	}

	// Yell is finished, stop it
	return animationYellMsg(false)
}

func (model Animation) yellStop() tea.Msg {
	model.pause()

	if model.Repeat == 0 {
		return animationStopMsg{}
	}

	return animationYellMsg(true)
}

func (model Animation) pause() {
	duration := (model.Duration / 2) + rand.IntN(model.Duration)
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

type animationStopMsg struct{}
type animationYellMsg bool
