package mascot

import (
	_ "embed"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"io"
)

func NewCmd() *cobra.Command {
	// Flags
	var repeat int

	// Command
	cmd := &cobra.Command{
		Use:    "duck",
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return run(cmd.OutOrStdout(), repeat)
		},
	}

	// Set flags
	cmd.Flags().IntVarP(&repeat, "repeat", "n", 1, "")

	return cmd
}

var (
	//go:embed assets/frame.txt
	frame string
	//go:embed assets/frame_yell.txt
	frameYell string
)

func run(out io.Writer, repeat int) error {
	renderer := lipgloss.NewRenderer(out)

	model, err := tea.NewProgram(
		animation{
			duration:  345,
			repeat:    repeat,
			style:     lipgloss.NewStyle().Renderer(renderer),
			frame:     &frame,
			frameYell: &frameYell,
		},
		tea.WithOutput(out),
		tea.WithAltScreen(),
	).Run()

	if err != nil {
		return err
	}

	return model.(animation).err
}
