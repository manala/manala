package update

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"

	"github.com/spf13/cobra"
)

func NewCommand(log *slog.Logger, api *api.API) *cobra.Command {
	// Flags
	var (
		repositoryURL, repositoryRef, recipeName string
		recursive                                bool
	)

	// Command
	command := &cobra.Command{
		Use:               "update [dir]",
		Aliases:           []string{"up"},
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Synchronize project(s)",
		Long: `Update (manala update) will synchronize project(s), based on repository's recipe
and related variables defined in manifest (.manala.yaml).

Example: manala update -> resulting in an update in a project dir (default to the
current directory)`,
		RunE: func(command *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			// Context
			ctx := command.Context()
			ctx = app.WithRepositoryURL(ctx, repositoryURL)
			ctx = app.WithRepositoryRef(ctx, repositoryRef)
			ctx = app.WithRecipeName(ctx, recipeName)

			return run(ctx, log, api, dir, recursive)
		},
	}

	// Set flags
	command.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	command.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	command.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")
	command.Flags().BoolVarP(&recursive, "recursive", "r", false, "set recursive mode")

	return command
}

func run(ctx context.Context, log *slog.Logger, api *api.API, dir string, recursive bool) error {
	// Get repository loader
	repositoryLoader := api.NewRepositoryLoader(ctx)

	// Get recipe loader
	recipeLoader := api.NewRecipeLoader(ctx)

	if recursive {
		// Get project loader
		projectLoader := api.NewProjectLoader(repositoryLoader, recipeLoader)

		// Recursively load projects
		log.Info("loading projects recursive…")

		return projectLoader.LoadRecursive(dir,
			func(project app.Project) error {
				// Sync project
				log.Info("syncing project…")

				return api.NewProjectSyncer().Sync(project)
			},
		)
	}

	// Get project loader
	projectLoader := api.NewProjectLoader(repositoryLoader, recipeLoader,
		api.WithProjectLoaderFrom(true),
	)

	// Load project
	log.Info("loading project…")

	project, err := projectLoader.Load(dir)
	if err != nil {
		return err
	}

	// Sync project
	return api.NewProjectSyncer().Sync(project)
}
