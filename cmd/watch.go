package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"manala/app"
	"manala/internal/config"
)

type WatchCmd struct{}

func (cmd *WatchCmd) Command(conf *config.Config, logger *log.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:     "watch [dir]",
		Aliases: []string{"Watch project"},
		Short:   "Watch project",
		Long: `Watch (manala watch) will watch project, and launch update on changes.

Example: manala watch -> resulting in a watch in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			// App
			_ = conf.BindPFlags(command.PersistentFlags())
			manala := app.New(conf, logger)

			// Command
			flags := config.New()
			_ = flags.BindPFlags(command.Flags())
			return manala.Watch(
				append(args, ".")[0],
				flags.GetString("recipe"),
				flags.GetBool("all"),
				flags.GetBool("notify"),
			)
		},
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("repository", "o", "", "with repository")

	// Flags
	flags := command.Flags()
	flags.StringP("recipe", "i", "", "with recipe")
	flags.BoolP("all", "a", false, "watch recipe too")
	flags.BoolP("notify", "n", false, "use system notifications")

	return command
}
