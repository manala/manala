package update

import (
	"context"
	"io"
	"path/filepath"

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

			return run(ctx, log, api, out, dir, recursive)
		},
	}

	// Set flags
	command.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	command.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	command.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")
	command.Flags().BoolVarP(&recursive, "recursive", "r", false, "set recursive mode")

	return command
}

func run(ctx context.Context, log *log.Log, api *api.API, out io.Writer, dir string, recursive bool) error {
	var (
		project app.Project
		err     error
	)

	// Api
	repositoryLoader := api.NewRepositoryLoader(ctx)
	recipeLoader := api.NewRecipeLoader(ctx)
	projectSyncer := api.NewProjectSyncer()

	if recursive {
		// Get project loader
		projectLoader := api.NewProjectLoader(repositoryLoader, recipeLoader)

		// Recursively load projects
		log.Info("loading projects recursive…")
		err = projectLoader.LoadRecursive(dir,
			func(project app.Project) error {
				// Sync project
				log.Info("syncing project…")
				err = projectSyncer.Sync(project)
				if err != nil {
					return err
				}

				return nil
			},
		)
		if err != nil {
			return err
		}

		return nil
	}

	// Get project loader
	projectLoader := api.NewProjectLoader(repositoryLoader, recipeLoader,
		api.WithProjectLoaderFrom(true),
	)

	// Load project
	log.Info("loading project…")
	project, err = projectLoader.Load(dir)
	if err != nil {
		return err
	}

	// Sync project
	log.Info("syncing project…")
	err = projectSyncer.Sync(project)
	if err != nil {
		return err
	}

	lipgloss.Fprintln(out, cmd.Styles.Primary.Render("project successfully updated"))

	return nil
}
