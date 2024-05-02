package watch

import (
	"github.com/spf13/cobra"
	"manala/app"
	"manala/app/api"
	"manala/internal/notifier"
	"manala/internal/ui"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func NewCMd(api *api.Api, out ui.Output, notifier notifier.Notifier) *cobra.Command {
	// Flags
	var repositoryUrl, repositoryRef, recipeName string
	var all, notify bool

	// Command
	cmd := &cobra.Command{
		Use:               "watch [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Watch project",
		Long: `Watch (manala watch) will watch project files, and launch update on changes.

Example: manala watch -> resulting in a watch in a project dir (default to the
current directory)`,
		RunE: func(_ *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			return run(api, out, notifier, dir, repositoryUrl, repositoryRef, recipeName, all, notify)
		},
	}

	// Set flags
	cmd.Flags().StringVarP(&repositoryUrl, "repository", "o", "", "use repository")
	cmd.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	cmd.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "watch recipe too")
	cmd.Flags().BoolVarP(&notify, "notify", "n", false, "use system notifications")

	return cmd
}

func run(api *api.Api, out ui.Output, notifier notifier.Notifier, dir, repositoryUrl, repositoryRef, recipeName string, all, notify bool) error {
	// Get repository loader
	repositoryLoader := api.NewRepositoryLoader(
		api.WithRepositoryLoaderUrl(repositoryUrl),
		api.WithRepositoryLoaderRef(repositoryRef),
	)

	// Get recipe loader
	recipeLoader := api.NewRecipeLoader(
		api.WithRecipeLoaderName(recipeName),
	)

	// Get project loader
	projectLoader := api.NewProjectLoader(repositoryLoader, recipeLoader,
		api.WithProjectLoaderFrom(true),
	)

	// Load project
	project, err := projectLoader.Load(dir)
	if err != nil {
		return err
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Watch project
	if err = api.NewProjectWatcher().Watch(project, all, func(project app.Project) app.Project {
		// Load project
		if project, err = projectLoader.Load(project.Dir()); err != nil {
			out.Error(err)
			if notify {
				notifier.Error(err)
			}
			return nil
		}

		// Sync project
		if err = api.NewProjectSyncer().Sync(project); err != nil {
			out.Error(err)
			if notify {
				notifier.Error(err)
			}
			return project
		}

		if notify {
			notifier.Message("Project synced")
		}

		return project
	}, done); err != nil {
		return err
	}

	return nil
}
