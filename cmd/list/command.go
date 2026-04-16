package list

import (
	"context"
	"io"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/cmd"
	"github.com/manala/manala/internal/log"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

func NewCommand(log *log.Log, api *api.API, out io.Writer) *cobra.Command {
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

func run(ctx context.Context, log *log.Log, api *api.API, out io.Writer) error {
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
		lipgloss.Fprintln(out, cmd.Styles.Primary.Render(recipe.Name()))
		lipgloss.Fprintln(out, "  "+cmd.Styles.Secondary.Render(recipe.Description()))
	}

	return nil
}
