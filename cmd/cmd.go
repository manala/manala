package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

func NewCmd(version string, in io.Reader, out io.Writer, err io.Writer) *cobra.Command {
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
	cmd.SetIn(in)
	cmd.SetOut(out)
	cmd.SetErr(err)

	return cmd
}
