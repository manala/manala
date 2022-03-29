package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/app"
)

type UpdateCmd struct{}

func (cmd *UpdateCmd) Command(config *viper.Viper, logger *log.Logger) *cobra.Command {
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
			recipe, _ := flags.GetString("recipe")
			recursive, _ := flags.GetBool("recursive")

			// Command
			return manala.Update(
				dir,
				recipe,
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
