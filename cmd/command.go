package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCommand(version string, in io.Reader, out, err io.Writer) *cobra.Command {
	// Command
	command := &cobra.Command{
		Use:               "manala",
		Version:           version,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Short:             "Let your project's plumbing up to date",
		Long: `Manala synchronize some boring parts of your projects, such as makefile targets,
virtualization and provisioning files...

Recipes are pulled from a git repository, or a local directory.`,
	}

	// Set streams
	command.SetIn(in)
	command.SetOut(out)
	command.SetErr(err)

	return command
}
