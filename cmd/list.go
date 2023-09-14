package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/app/config"
	"manala/internal/ui"
)

func newListCmd(config config.Config, log *slog.Logger, out ui.Output) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"ls"},
		Args:              cobra.NoArgs,
		DisableAutoGenTag: true,
		Short:             "List recipes",
		Long: `List (manala list) will list recipes available on
repository.

Example: manala list -> resulting in a recipes list display`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Api options
			var apiOptions []api.Option

			// Flag - Repository url
			if cmd.Flags().Changed("repository") {
				repositoryUrl, _ := cmd.Flags().GetString("repository")
				apiOptions = append(apiOptions, api.WithRepositoryUrl(repositoryUrl))
			}

			// Flag - Repository ref
			if cmd.Flags().Changed("ref") {
				repositoryRef, _ := cmd.Flags().GetString("ref")
				apiOptions = append(apiOptions, api.WithRepositoryRef(repositoryRef))
			}

			// Api
			api := api.New(
				config,
				log,
				out,
				apiOptions...,
			)

			// Load preceding repository
			repository, err := api.LoadPrecedingRepository()
			if err != nil {
				return err
			}

			// List
			var recipes []app.Recipe
			if err := api.WalkRepositoryRecipes(repository, func(recipe app.Recipe) error {
				recipes = append(recipes, recipe)
				return nil
			}); err != nil {
				return err
			}

			return out.List(
				fmt.Sprintf("Recipes available in %s", repository.Url()),
				api.NewUiRecipeList(recipes),
			)
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().String("ref", "", "use repository ref")

	return cmd
}
