package watch

import (
	"context"
	"io"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/cmd"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/notify"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

func NewCommand(log *log.Log, api *api.API, out io.Writer, notifier *notify.Notifier) *cobra.Command {
	// Flags
	var (
		repositoryURL, repositoryRef, recipeName string
		all, notify                              bool
	)

	// Command
	command := &cobra.Command{
		Use:               "watch [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Watch project",
		Long: `Watch (manala watch) will watch project files, and launch update on changes.

Example: manala watch -> resulting in a watch in a project dir (default to the
current directory)`,
		RunE: func(command *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			// Context
			ctx := command.Context()
			ctx = app.WithRepositoryURL(ctx, repositoryURL)
			ctx = app.WithRepositoryRef(ctx, repositoryRef)
			ctx = app.WithRecipeName(ctx, recipeName)

			return run(ctx, log, api, out, notifier, dir, all, notify)
		},
	}

	// Set flags
	command.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	command.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	command.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")
	command.Flags().BoolVarP(&all, "all", "a", false, "watch recipe too")
	command.Flags().BoolVarP(&notify, "notify", "n", false, "use system notifications")

	return command
}

func run(ctx context.Context, log *log.Log, api *api.API, out io.Writer, notifier *notify.Notifier, dir string, all, notify bool) error {
	var (
		project app.Project
		err     error
	)
	// Api
	repositoryLoader := api.NewRepositoryLoader(ctx)
	recipeLoader := api.NewRecipeLoader(ctx)
	projectLoader := api.NewProjectLoader(repositoryLoader, recipeLoader,
		api.WithProjectLoaderFrom(true),
	)
	projectSyncer := api.NewProjectSyncer()

	// Load project
	log.Info("loading project…")
	project, err = projectLoader.Load(dir)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Watch project
	log.Info("watching project…")
	err = NewWatcher(log, all).
		Watch(ctx, project, func(project app.Project) app.Project {
			// Load project
			log.Info("loading project…")
			if project, err = projectLoader.Load(project.Dir()); err != nil {
				log.Error(err)

				if notify {
					notifier.Error(err)
				}

				return nil
			}

			// Sync project
			log.Info("syncing project…")
			if err = projectSyncer.Sync(project); err != nil {
				log.Error(err)

				if notify {
					notifier.Error(err)
				}
			} else {
				notifier.Message("Project synced")
			}

			lipgloss.Fprintln(out, cmd.Styles.Primary.Render("project successfully updated"))

			return project
		})
	if err != nil {
		return err
	}

	return nil
}
