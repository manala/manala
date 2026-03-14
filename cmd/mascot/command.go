package mascot

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCommand(in io.Reader, out io.Writer) *cobra.Command {
	// Flags
	var repeat int

	// Command
	command := &cobra.Command{
		Use:    "goat",
		Hidden: true,
		RunE: func(command *cobra.Command, _ []string) error {
			return run(in, out, repeat)
		},
	}

	// Set flags
	command.Flags().IntVarP(&repeat, "repeat", "n", 1, "")

	return command
}

func run(in io.Reader, out io.Writer, repeat int) error {
	// Run mascot
	err := RunMascot(in, out, repeat)
	if err != nil {
		return err
	}

	return nil
}
