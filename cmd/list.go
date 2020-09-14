package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/loaders"
	"manala/models"
)

// ListCmd represents the list command
func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List recipes",
		Long: `List (manala list) will list recipes available on
repository.

Example: manala list -> resulting in a recipes list display`,
		RunE: listRun,
		Args: cobra.NoArgs,
	}

	addRepositoryFlag(cmd, "use repository")

	return cmd
}

func listRun(cmd *cobra.Command, args []string) error {
	// Loaders
	repoLoader := loaders.NewRepositoryLoader(
		viper.GetString("cache_dir"),
		viper.GetString("repository"),
	)
	recLoader := loaders.NewRecipeLoader()

	// Load repository
	repoName, _ := cmd.Flags().GetString("repository")
	repo, err := repoLoader.Load(repoName)
	if err != nil {
		return err
	}

	// Walk into recipes
	if err := recLoader.Walk(repo, func(rec models.RecipeInterface) {
		cmd.Printf("%s: %s\n", rec.Name(), rec.Description())
	}); err != nil {
		return err
	}

	return nil
}
