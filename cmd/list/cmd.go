package list

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/internal/ui"
)

func NewCmd(log *slog.Logger, api *api.Api, output ui.Output) *cobra.Command {
	// Flags
	var repositoryUrl, repositoryRef string

	// Command
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"ls"},
		Args:              cobra.NoArgs,
		DisableAutoGenTag: true,
		Short:             "List recipes",
		Long: `List (manala list) will list recipes available on repository.

Example: manala list -> resulting in a recipes list display`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Context
			ctx := cmd.Context()
			ctx = app.WithRepositoryUrl(ctx, repositoryUrl)
			ctx = app.WithRepositoryRef(ctx, repositoryRef)

			return run(ctx, log, api, output)
		},
	}

	// Set flags
	cmd.Flags().StringVarP(&repositoryUrl, "repository", "o", "", "use repository")
	cmd.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")

	return cmd
}

func run(ctx context.Context, log *slog.Logger, api *api.Api, output ui.Output) error {
	// Get repository loader
	repositoryLoader := api.NewRepositoryLoader(ctx)

	// Load repository
	log.Info("loading repository…")
	repository, err := repositoryLoader.Load("")
	if err != nil {
		return err
	}

	// Get recipe loader
	recipeLoader := api.NewRecipeLoader(ctx)

	// Load recipes
	log.Info("loading recipes…")
	recipes, err := recipeLoader.LoadAll(repository)
	if err != nil {
		return err
	}

	return output.List(
		fmt.Sprintf("Recipes available in %s", repository.Url()),
		NewUiRecipeList(recipes),
	)
}
