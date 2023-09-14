package charm

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"manala/internal/ui/components"
	"math/rand"
	"time"
)

func (adapter *Adapter) Animate(animation components.Animation, repeat int) error {
	renderer := adapter.outRenderer

	model := newAnimationModel(animation, repeat, renderer)

	_model, err := tea.NewProgram(
		model,
		tea.WithOutput(renderer.Output()),
		tea.WithAltScreen(),
	).Run()

	if err != nil {
		return err
	}

	return _model.(animationModel).err
}

/*********/
/* Model */
/*********/

func newAnimationModel(animation components.Animation, repeat int, renderer *lipgloss.Renderer) animationModel {
	return animationModel{
		animation: animation,
		repeat:    repeat,
		style:     animationStyle.New(renderer),
	}
}

type animationModel struct {
	animation components.Animation
	repeat    int
	yell      bool
	width     int
	height    int
	err       error
	style     *style
}

func (model animationModel) Init() tea.Cmd {
	return func() tea.Msg {
		if model.repeat == 0 {
			return animationStopMsg{}
		}

		return animationYellStartMsg{}
	}
}

func (model animationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch _msg := msg.(type) {
	// Size
	case tea.WindowSizeMsg:
		model.width, model.height = _msg.Width, _msg.Height
	// Error
	case error:
		model.err = _msg
		cmd = tea.Quit
	// Keys
	case tea.KeyMsg:
		switch {
		// Keys - Quit
		case key.Matches(_msg, QuitOverallKeys):
			cmd = tea.Quit
		}
	// Animation - Stop
	case animationStopMsg:
		cmd = tea.Quit
	// Animation - Yell
	case animationYellStartMsg:
		model.yell = true
		cmd = model.yellStart
	case animationYellStopMsg:
		model.yell = false
		if model.repeat > 0 {
			model.repeat--
		}
		cmd = model.yellStop
	}

	return model, cmd
}

func (model animationModel) View() string {
	frame := model.animation.Frame
	if model.yell {
		frame = model.animation.FrameYell
	}

	// Render
	return lipgloss.Place(
		model.width, model.height,
		lipgloss.Center, lipgloss.Center,
		model.style.Style().
			MaxWidth(model.width).
			MaxHeight(model.height).
			Render(*frame),
	)
}

func (model animationModel) yellStart() tea.Msg {
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
	return animationYellStopMsg{}
}

func (model animationModel) yellStop() tea.Msg {
	model.pause()

	if model.repeat == 0 {
		return animationStopMsg{}
	}

	return animationYellStartMsg{}
}

func (model animationModel) pause() {
	duration := (model.animation.Duration / 2) + rand.Intn(model.animation.Duration)
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

type animationStopMsg struct{}
type animationYellStartMsg struct{}
type animationYellStopMsg struct{}
