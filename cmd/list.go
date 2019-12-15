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
	repo, err := repository.Load(viper.GetString("repository"), viper.GetString("cache_dir"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Walk into recipes
	err = recipe.Walk(repo, func(rec recipe.Interface) {
		fmt.Printf("%s: %s\n", rec.GetName(), rec.GetConfig().Description)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
