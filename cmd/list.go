package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"manala/app"
	"manala/internal/config"
)

type ListCmd struct{}

func (cmd *ListCmd) Command(conf *config.Config, logger *log.Logger) *cobra.Command {
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
			_ = conf.BindPFlags(command.PersistentFlags())
			manala := app.New(conf, logger)

			// Command
			recipes, err := manala.List()
			if err != nil {
				return err
			}

			for _, recipe := range recipes {
				_, _ = fmt.Printf("%s: %s\n", recipe.Name(), recipe.Description())
			}

			return nil
		},
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("repository", "o", "", "use repository")

	return command
}
