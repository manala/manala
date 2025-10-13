package watch

import (
	"context"
	"log/slog"
	"os/signal"
	"path/filepath"
	"syscall"

	"manala/app"
	"manala/app/api"
	"manala/internal/notifier"
	"manala/internal/ui"

	"github.com/spf13/cobra"
)

func NewCmd(log *slog.Logger, api *api.API, output ui.Output, notifier notifier.Notifier) *cobra.Command {
	// Flags
	var (
		repositoryURL, repositoryRef, recipeName string
		all, notify                              bool
	)

	// Command
	cmd := &cobra.Command{
		Use:               "watch [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Watch project",
		Long: `Watch (manala watch) will watch project files, and launch update on changes.

Example: manala watch -> resulting in a watch in a project dir (default to the
current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			// Context
			ctx := cmd.Context()
			ctx = app.WithRepositoryURL(ctx, repositoryURL)
			ctx = app.WithRepositoryRef(ctx, repositoryRef)
			ctx = app.WithRecipeName(ctx, recipeName)

			return run(ctx, log, api, output, notifier, dir, all, notify)
		},
	}

	// Set flags
	cmd.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	cmd.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	cmd.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "watch recipe too")
	cmd.Flags().BoolVarP(&notify, "notify", "n", false, "use system notifications")

	return cmd
}

func run(ctx context.Context, log *slog.Logger, api *api.API, output ui.Output, notifier notifier.Notifier, dir string, all, notify bool) error {
	// Get repository loader
	repositoryLoader := api.NewRepositoryLoader(ctx)

	// Get recipe loader
	recipeLoader := api.NewRecipeLoader(ctx)

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

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Watch project
	log.Info("watching project…")

	return NewWatcher(log, all).
		Watch(ctx, project, func(project app.Project) app.Project {
			// Load project
			log.Info("loading project…")

			if project, err = projectLoader.Load(project.Dir()); err != nil {
				output.Error(err)

				if notify {
					notifier.Error(err)
				}

				return nil
			}

			// Sync project
			log.Info("syncing project…")

			if err = api.NewProjectSyncer().Sync(project); err != nil {
				output.Error(err)

				if notify {
					notifier.Error(err)
				}
			} else {
				notifier.Message("Project synced")
			}

			return project
		})
}
