package cmd

import (
	"github.com/spf13/cobra"
	"manala/config"
)

type RootCmd struct {
	Conf config.Config
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
		Version:           cmd.Conf.Version(),
		DisableAutoGenTag: true,
	}

	pFlags := command.PersistentFlags()

	// Cache dir
	pFlags.StringP("cache-dir", "c", "", "use cache directory")
	cmd.Conf.BindCacheDirFlag(pFlags.Lookup("cache-dir"))

	// Debug
	pFlags.BoolP("debug", "d", false, "set debug mode")
	cmd.Conf.BindDebugFlag(pFlags.Lookup("debug"))

	return command
}
