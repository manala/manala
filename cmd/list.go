package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"manala/app"
)

type ListCmd struct {
	App *app.App
	Out io.Writer
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
			// App
			recipes, err := cmd.App.List()
			if err != nil {
				return err
			}

			for _, recipe := range recipes {
				_, _ = fmt.Fprintf(cmd.Out, "%s: %s\n", recipe.Name(), recipe.Description())
			}

			return nil
		},
	}

	flags := command.Flags()

	// Repository
	flags.StringP("repository", "o", "", "use repository source")
	cmd.App.Config.BindPFlag("repository", flags.Lookup("repository"))

	return command
}
