package cmd

import (
	"github.com/spf13/cobra"
	"manala/app"
)

type WatchCmd struct {
	App *app.App
}

func (cmd *WatchCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:     "watch [dir]",
		Aliases: []string{"Watch project"},
		Short:   "Watch project",
		Long: `Watch (manala watch) will watch project, and launch update on changes.

Example: manala watch -> resulting in a watch in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			// Get directory from first command arg
			dir := "."
			if len(args) != 0 {
				dir = args[0]
			}

			flags := command.Flags()

			withRecipeName, _ := flags.GetString("recipe")

			watchAll, _ := flags.GetBool("all")
			useNotify, _ := flags.GetBool("notify")

			// App
			return cmd.App.Watch(
				dir,
				withRecipeName,
				watchAll,
				useNotify,
			)
		},
	}

	flags := command.Flags()

	// Repository
	flags.StringP("repository", "o", "", "with repository source")
	cmd.App.Config.BindPFlag("repository", flags.Lookup("repository"))

	// Recipe
	flags.StringP("recipe", "i", "", "with recipe name")

	flags.BoolP("all", "a", false, "watch recipe too")
	flags.BoolP("notify", "n", false, "use system notifications")

	return command
}
