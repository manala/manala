package cmd

import (
	"github.com/spf13/cobra"
	"manala/app"
)

type UpdateCmd struct {
	App *app.App
}

func (cmd *UpdateCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:     "update [dir]",
		Aliases: []string{"up"},
		Short:   "Update project",
		Long: `Update (manala update) will update project, based on
recipe and related variables defined in manala.yaml.

Example: manala update -> resulting in an update in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			// Config
			cmd.App.Config.BindPFlags(command.PersistentFlags())

			// Get directory from first command arg
			dir := "."
			if len(args) != 0 {
				dir = args[0]
			}

			// Flags
			flags := command.Flags()
			withRecipeName, _ := flags.GetString("recipe")
			recursive, _ := flags.GetBool("recursive")

			// App
			return cmd.App.Update(
				dir,
				withRecipeName,
				recursive,
			)
		},
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("repository", "o", "", "with repository source")

	// Flags
	flags := command.Flags()
	flags.StringP("recipe", "i", "", "with recipe name")
	flags.BoolP("recursive", "r", false, "set recursive mode")

	return command
}
