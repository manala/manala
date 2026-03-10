package mascot

import (
	_ "embed"
	"io"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	// Flags
	var repeat int

	// Command
	command := &cobra.Command{
		Use:    "duck",
		Hidden: true,
		RunE: func(command *cobra.Command, _ []string) error {
			return run(command.InOrStdin(), command.OutOrStdout(), repeat)
		},
	}

	// Set flags
	command.Flags().IntVarP(&repeat, "repeat", "n", 1, "")

	return command
}

var (
	//go:embed assets/frame.txt
	frame string
	//go:embed assets/frame_yell.txt
	frameYell string
	//go:embed assets/audio_yell.wav
	audioYell []byte
)

func run(in io.Reader, out io.Writer, repeat int) error {
	renderer := lipgloss.NewRenderer(out)

	model, err := tea.NewProgram(
		Animation{
			Title:    "Quack Quack",
			Duration: 345,
			Repeat:   repeat,
			Style:    lipgloss.NewStyle().Renderer(renderer),
			QuitKey: key.NewBinding(
				key.WithKeys("ctrl+c", "esc", "q"),
			),
			Frame:     &frame,
			FrameYell: &frameYell,
			AudioYell: &audioYell,
		},
		tea.WithInput(in),
		tea.WithOutput(out),
		tea.WithAltScreen(),
	).Run()
	if err != nil {
		return err
	}

	return model.(Animation).Err
}
