package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/app"
)

type WatchCmd struct{}

func (cmd *WatchCmd) Command(config *viper.Viper, logger *log.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:     "watch [dir]",
		Aliases: []string{"Watch project"},
		Short:   "Watch project",
		Long: `Watch (manala watch) will watch project, and launch update on changes.

Example: manala watch -> resulting in a watch in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			// Config
			_ = config.BindPFlags(command.PersistentFlags())

			// App
			manala := app.New(
				app.WithConfig(config),
				app.WithLogger(logger),
			)

			// Get directory from first command arg
			dir := "."
			if len(args) != 0 {
				dir = args[0]
			}

			// Flags
			flags := command.Flags()
			withRecipeName, _ := flags.GetString("recipe")
			watchAll, _ := flags.GetBool("all")
			useNotify, _ := flags.GetBool("notify")

			// Command
			return manala.Watch(
				dir,
				withRecipeName,
				watchAll,
				useNotify,
			)
		},
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("repository", "o", "", "with repository source")

	// Flags
	flags := command.Flags()
	flags.StringP("recipe", "i", "", "with recipe name")
	flags.BoolP("all", "a", false, "watch recipe too")
	flags.BoolP("notify", "n", false, "use system notifications")

	return command
}
