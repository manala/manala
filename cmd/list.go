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
			// Config
			cmd.App.Config.BindPFlags(command.PersistentFlags())

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

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("repository", "o", "", "use repository source")

	return command
}
