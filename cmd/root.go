package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootCmd struct{}

func (cmd *RootCmd) Command(config *viper.Viper, logger *log.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "manala",
		Short: "Let your project's plumbing up to date",
		Long: `Manala synchronize some boring parts of your projects,
such as makefile targets, virtualization and provisioning files...

Recipes are pulled from a git repository, or a local directory.`,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Version:           config.GetString("version"),
		DisableAutoGenTag: true,
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("cache-dir", "c", "", "use cache directory")
	pFlags.BoolP("debug", "d", false, "set debug mode")

	_ = config.BindPFlags(pFlags)

	// Debug
	cobra.OnInitialize(func() {
		if config.GetBool("debug") {
			logger.Level = log.DebugLevel
		}
	})

	return command
}
