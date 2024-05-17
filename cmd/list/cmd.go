package list

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app/api"
	"manala/internal/ui"
)

func NewCmd(log *slog.Logger, api *api.Api, out ui.Output) *cobra.Command {
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
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(log, api, out, repositoryUrl, repositoryRef)
		},
	}

	// Set flags
	cmd.Flags().StringVarP(&repositoryUrl, "repository", "o", "", "use repository")
	cmd.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")

	return cmd
}

func run(log *slog.Logger, api *api.Api, out ui.Output, repositoryUrl, repositoryRef string) error {
	// Get repository loader
	repositoryLoader := api.NewRepositoryLoader(
		api.WithRepositoryLoaderRef(repositoryRef),
	)

	// Load repository
	log.Info("loading repository…")
	repository, err := repositoryLoader.Load(repositoryUrl)
	if err != nil {
		return err
	}

	// Get recipe loader
	recipeLoader := api.NewRecipeLoader()

	// Load recipes
	log.Info("loading recipes…")
	recipes, err := recipeLoader.LoadAll(repository)
	if err != nil {
		return err
	}

	return out.List(
		fmt.Sprintf("Recipes available in %s", repository.Url()),
		NewUiRecipeList(recipes),
	)
}
