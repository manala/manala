package list

import (
	"context"
	"log/slog"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/internal/ui"

	"github.com/spf13/cobra"
)

func NewCommand(log *slog.Logger, api *api.API, output ui.Output) *cobra.Command {
	// Flags
	var repositoryURL, repositoryRef string

	// Command
	command := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"ls"},
		Args:              cobra.NoArgs,
		DisableAutoGenTag: true,
		Short:             "List recipes",
		Long: `List (manala list) will list recipes available on repository.

Example: manala list -> resulting in a recipes list display`,
		RunE: func(command *cobra.Command, _ []string) error {
			// Context
			ctx := command.Context()
			ctx = app.WithRepositoryURL(ctx, repositoryURL)
			ctx = app.WithRepositoryRef(ctx, repositoryRef)

			return run(ctx, log, api, output)
		},
	}

	// Set flags
	command.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	command.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")

	return command
}

func run(ctx context.Context, log *slog.Logger, api *api.API, output ui.Output) error {
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
		"Recipes available in "+repository.URL(),
		NewUIRecipeList(recipes),
	)
}
