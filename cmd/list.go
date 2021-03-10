package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"manala/loaders"
	"manala/models"
)

type ListCmd struct {
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

			repoSrc, _ := flags.GetString("repository")

			return cmd.Run(repoSrc)
		},
	}

	flags := command.Flags()

	flags.StringP("repository", "o", "", "use repository source")

	return command
}

func (cmd *ListCmd) Run(repoSrc string) error {
	// Load repository
	repo, err := cmd.RepositoryLoader.Load(repoSrc)
	if err != nil {
		return err
	}

	// Walk into recipes
	if err := cmd.RecipeLoader.Walk(repo, func(rec models.RecipeInterface) {
		fmt.Fprintf(cmd.Out, "%s: %s\n", rec.Name(), rec.Description())
	}); err != nil {
		return err
	}

	return nil
}
