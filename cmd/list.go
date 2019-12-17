package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/pkg/recipe"
	"manala/pkg/repository"
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
	// Load repository
	repo := repository.New(viper.GetString("repository"))
	if err := repo.Load(viper.GetString("cache_dir")); err != nil {
		log.Fatal(err.Error())
	}

	// Walk into recipes
	if err := repo.WalkRecipes(func(rec recipe.Interface) {
		fmt.Printf("%s: %s\n", rec.GetName(), rec.GetConfig().Description)
	}); err != nil {
		log.Fatal(err.Error())
	}
}
