package cmd

import (
	_ "embed"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

func newMascotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "duck",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			repeat, _ := cmd.Flags().GetInt("repeat")

			// Model
			model := &mascotModel{
				duration: 345,
				style:    lipgloss.NewStyle(),
				repeat:   repeat,
			}

			// Program
			if err := tea.NewProgram(
				model,
				tea.WithAltScreen(),
				tea.WithOutput(cmd.OutOrStdout()),
			).Start(); err != nil {
				return err
			}

			return model.err
		},
	}

	// Flags
	cmd.Flags().IntP("repeat", "n", 1, "")

	return cmd
}

//go:embed resources/mascot.txt
var mascotText string

//go:embed resources/mascot_yell.txt
var mascotTextYell string

type mascotModel struct {
	duration     int
	windowWidth  int
	windowHeight int
	style        lipgloss.Style
	yell         bool
	repeat       int
	err          error
}

func (model *mascotModel) Init() tea.Cmd {
	return func() tea.Msg {
		if model.repeat == 0 {
			return "stop"
		}

		return "yellStart"
	}
}

func (model *mascotModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window size
	case tea.WindowSizeMsg:
		model.windowWidth, model.windowHeight = msg.Width, msg.Height
		return model, nil

	// Keys
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return model, tea.Quit
		}

	// Error
	case error:
		model.err = msg
		return model, tea.Quit

	// Yell / Stop
	case string:
		switch msg {
		case "yellStart":
			model.yell = true
			return model, model.yellStart
		case "yellStop":
			model.yell = false
			return model, model.yellStop
		case "stop":
			return model, tea.Quit
		}
	}

	return model, nil
}

func (model *mascotModel) View() string {
	// Wait for model to be ready
	if model.windowWidth == 0 && model.windowHeight == 0 {
		return ""
	}

	if model.yell {
		model.style = model.style.SetString(mascotTextYell)
	} else {
		model.style = model.style.SetString(mascotText)
	}

	// Render
	return lipgloss.Place(model.windowWidth, model.windowHeight, lipgloss.Center, lipgloss.Center,
		model.style.
			MaxHeight(model.windowHeight).
			String(),
	)
}

func (model *mascotModel) yellStart() tea.Msg {
	// Has yell audio ?
	var modelInterface interface{} = model
	if _model, ok := modelInterface.(interface{ yellAudio() error }); ok {
		if err := _model.yellAudio(); err != nil {
			return err
		}
	} else {
		model.pause()
	}

	// Yell is finished, stop it
	return "yellStop"
}

func (model *mascotModel) yellStop() tea.Msg {
	model.pause()

	// Yell is stopped, do we need to start it again ?
	if model.repeat < 0 {
		return "yellStart"
	} else if model.repeat--; model.repeat == 0 {
		return "stop"
	}

	return "yellStart"
}

func (model *mascotModel) pause() {
	duration := (model.duration / 2) + rand.Intn(model.duration)
	time.Sleep(time.Duration(duration) * time.Millisecond)
}
