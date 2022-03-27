package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"manala/app"
	"manala/config"
)

type ListCmd struct {
	App  *app.App
	Conf config.Config
	Out  io.Writer
}

func (cmd *ListCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List recipes",
		Long: `List (manala list) will list recipes available on
repository.

Example: manala list -> resulting in a recipes list display`,
		Args:              cobra.NoArgs,
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			flags := command.Flags()

			cmd.Conf.BindRepositoryFlag(flags.Lookup("repository"))

			// App
			return cmd.App.List(
				cmd.Conf.Repository(),
				cmd.Out,
			)
		},
	}

	flags := command.Flags()

	flags.StringP("repository", "o", "", "use repository source")

	return command
}
