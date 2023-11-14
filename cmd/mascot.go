package cmd

import (
	_ "embed"
	"github.com/spf13/cobra"
	"manala/internal/ui"
	"manala/internal/ui/components"
)

var (
	//go:embed resources/mascot.txt
	mascotText string
	//go:embed resources/mascot_yell.txt
	mascotTextYell string
	//go:embed resources/mascot_yell.wav
	mascotAudioYell []byte
)

func newMascotCmd(out ui.Output) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "duck",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			repeat, _ := cmd.Flags().GetInt("repeat")

			return out.Animate(
				components.Animation{
					Frame:     &mascotText,
					FrameYell: &mascotTextYell,
					AudioYell: &mascotAudioYell,
					Duration:  345,
				},
				repeat,
			)
		},
	}

	// Flags
	cmd.Flags().IntP("repeat", "n", 1, "")

	return cmd
}
