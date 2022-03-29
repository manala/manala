package cmd

import (
	"github.com/spf13/cobra"
	"manala/app"
)

type RootCmd struct {
	App          *app.App
	OnInitialize func()
}

func (cmd *RootCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "manala",
		Short: "Let your project's plumbing up to date",
		Long: `Manala synchronize some boring parts of your projects,
such as makefile targets, virtualization and provisioning files...

Recipes are pulled from a git repository, or a local directory.`,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Version:           cmd.App.Config.GetString("version"),
		DisableAutoGenTag: true,
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("cache-dir", "c", "", "use cache directory")
	pFlags.BoolP("debug", "d", false, "set debug mode")

	cmd.App.Config.BindPFlags(pFlags)

	// Initialize
	cobra.OnInitialize(cmd.OnInitialize)

	return command
}
