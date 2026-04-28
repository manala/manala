package list

import (
	"context"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/output"

	"github.com/spf13/cobra"
)

func NewCommand(log *log.Log, api *api.API, out output.Output) *cobra.Command {
	// Flags
	var (
		repositoryURL string
		repositoryRef string
	)

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

			return run(ctx, log, api, out)
		},
	}

	// Set flags
	command.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	command.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")

	return command
}

func run(ctx context.Context, log *log.Log, api *api.API, out output.Output) error {
	var (
		repository app.Repository
		recipes    []app.Recipe
		err        error
	)

	// Api
	repositoryLoader := api.NewRepositoryLoader(ctx)
	recipeLoader := api.NewRecipeLoader(ctx)

	// Load repository
	log.Info("loading repository…")
	repository, err = repositoryLoader.Load("")
	if err != nil {
		return err
	}

	// Load recipes
	log.Info("loading recipes…")
	recipes, err = recipeLoader.LoadAll(repository)
	if err != nil {
		return err
	}

	for _, recipe := range recipes {
		out.Println(out.Style().Render(recipe.Name()))
		out.Println("  " + out.MutedStyle().Render(recipe.Description()))
	}

	return nil
}
