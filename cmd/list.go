package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"manala/config"
	"manala/loaders"
	"manala/models"
)

type ListCmd struct {
	Conf             *config.Config
	RepositoryLoader loaders.RepositoryLoaderInterface
	RecipeLoader     loaders.RecipeLoaderInterface
	Out              io.Writer
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

			return cmd.Run()
		},
	}

	flags := command.Flags()

	flags.StringP("repository", "o", "", "use repository source")

	return command
}

func (cmd *ListCmd) Run() error {
	// Load repository
	repo, err := cmd.RepositoryLoader.Load(cmd.Conf.Repository())
	if err != nil {
		return err
	}

	// Walk into recipes
	if err := cmd.RecipeLoader.Walk(repo, func(rec models.RecipeInterface) {
		_, _ = fmt.Fprintf(cmd.Out, "%s: %s\n", rec.Name(), rec.Description())
	}); err != nil {
		return err
	}

	return nil
}
