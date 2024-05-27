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

func (model animation) Init() tea.Cmd {
	return tea.Sequence(
		tea.SetWindowTitle("Quack Quack"),
		func() tea.Msg {
			if model.repeat == 0 {
				return animationStopMsg{}
			}

			return animationYellMsg(true)
		},
	)
}

func (model animation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Size
	case tea.WindowSizeMsg:
		model.width, model.height = msg.Width, msg.Height
	// Error
	case error:
		model.err = msg
		cmd = tea.Quit
	// Keys
	case tea.KeyMsg:
		switch msg.String() {
		// Keys - Quit
		case "ctrl+c", "esc", "q":
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
			if model.repeat > 0 {
				model.repeat--
			}
			cmd = model.yellStop
		}
	}

	return model, cmd
}

func (model animation) View() string {
	frame := model.frame
	if model.yell {
		frame = model.frameYell
	}

	// Render
	return lipgloss.Place(
		model.width, model.height,
		lipgloss.Center, lipgloss.Center,
		model.style.
			MaxWidth(model.width).
			MaxHeight(model.height).
			Render(*frame),
	)
}

func (model animation) yellStart() tea.Msg {
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

func (model animation) yellStop() tea.Msg {
	model.pause()

	if model.repeat == 0 {
		return animationStopMsg{}
	}

	return animationYellMsg(true)
}

func (model animation) pause() {
	duration := (model.duration / 2) + rand.IntN(model.duration)
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

type animationStopMsg struct{}
type animationYellMsg bool
