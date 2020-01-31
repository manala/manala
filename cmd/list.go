package cmd

import (
	"fmt"
	"github.com/apex/log"
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
		Run:  listRun,
		Args: cobra.NoArgs,
	}

	return cmd
}

func listRun(cmd *cobra.Command, args []string) {
	// Loaders
	repoLoader := loaders.NewRepositoryLoader(viper.GetString("cache_dir"))
	recLoader := loaders.NewRecipeLoader()

	// Load repository
	repo, err := repoLoader.Load(viper.GetString("repository"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Walk into recipes
	if err := recLoader.Walk(repo, func(rec models.RecipeInterface) {
		fmt.Printf("%s: %s\n", rec.Name(), rec.Description())
	}); err != nil {
		log.Fatal(err.Error())
	}
}
