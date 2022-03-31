package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"manala/app"
	"manala/internal/config"
)

type UpdateCmd struct{}

func (cmd *UpdateCmd) Command(conf *config.Config, logger *log.Logger) *cobra.Command {
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
			// App
			_ = conf.BindPFlags(command.PersistentFlags())
			manala := app.New(conf, logger)

			// Command
			flags := config.New()
			_ = flags.BindPFlags(command.Flags())
			return manala.Update(
				append(args, ".")[0],
				flags.GetString("recipe"),
				flags.GetBool("recursive"),
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
