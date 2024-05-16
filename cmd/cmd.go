package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

func NewCmd(version string, stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	// Command
	cmd := &cobra.Command{
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
	cmd.SetOut(stdOut)
	cmd.SetErr(stdErr)

	return cmd
}
