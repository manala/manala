package init

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
		repositoryURL string
		repositoryRef string
		recipeName    string
	)

	// Command
	command := &cobra.Command{
		Use:               "init [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a dir (default to the
current directory)`,
		RunE: func(command *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			// Context
			ctx := command.Context()
			ctx = app.WithRepositoryURL(ctx, repositoryURL)
			ctx = app.WithRepositoryRef(ctx, repositoryRef)
			ctx = app.WithRecipeName(ctx, recipeName)

			return run(ctx, log, api, out, dir)
		},
	}

	// Set flags
	command.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	command.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	command.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")

	return command
}

func run(ctx context.Context, log *log.Log, api *api.API, out io.Writer, dir string) error {
	var (
		repository    app.Repository
		dialogVariant DialogVariant
		project       app.Project
		err           error
	)

	// Api
	projectFinder := api.NewProjectFinder()
	repositoryLoader := api.NewRepositoryLoader(ctx)
	recipeLoader := api.NewRecipeLoader(ctx)
	projectCreator := api.NewProjectCreator()
	projectSyncer := api.NewProjectSyncer()

	// Check already existing project
	log.Info("finding project…")
	if projectFinder.Find(dir) {
		return &app.AlreadyExistingProjectError{Dir: dir}
	}

	// Load repository
	log.Info("loading repository…")
	repository, err = repositoryLoader.Load("")
	if err != nil {
		return err
	}

	if _, ok := app.RecipeName(ctx); ok {
		// Load recipe by context
		log.Info("loading recipe…")
		recipe, err := recipeLoader.Load(repository, "")
		if err != nil {
			return err
		}
		dialogVariant = DialogSingleVariant{Recipe: recipe}
	} else {
		// Select recipe
		log.Info("loading recipes…")
		recipes, err := recipeLoader.LoadAll(repository)
		if err != nil {
			return err
		}
		dialogVariant = DialogMultiVariant{Recipes: recipes}
	}

	// Run dialog
	outcome, err := RunDialog("Manala", dialogVariant)
	if err != nil {
		return err
	}

	// Create project
	log.Info("creating project…")
	project, err = projectCreator.Create(dir, outcome.Recipe, outcome.Vars)
	if err != nil {
		return err
	}

	// Sync project
	log.Info("syncing project…")
	err = projectSyncer.Sync(project)
	if err != nil {
		return err
	}

	lipgloss.Fprintln(out, cmd.Styles.Primary.Render("project successfully initialized"))

	return nil
}
