package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"manala/internal/config"
)

type RootCmd struct{}

func (cmd *RootCmd) Command(conf *config.Config, logger *log.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "manala",
		Short: "Let your project's plumbing up to date",
		Long: `Manala synchronize some boring parts of your projects,
such as makefile targets, virtualization and provisioning files...

Recipes are pulled from a git repository, or a local directory.`,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Version:           conf.GetString("version"),
		DisableAutoGenTag: true,
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("cache-dir", "c", "", "use cache directory")
	pFlags.BoolP("debug", "d", false, "set debug mode")

	_ = conf.BindPFlags(pFlags)

	// Debug
	cobra.OnInitialize(func() {
		if conf.GetBool("debug") {
			logger.Level = log.DebugLevel
		}
	})

	return command
}
