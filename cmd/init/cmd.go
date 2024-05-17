package init

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/internal/ui"
	"path/filepath"
)

func NewCmd(log *slog.Logger, api *api.Api, in ui.Input) *cobra.Command {
	// Flags
	var repositoryUrl, repositoryRef, recipeName string

	// Command
	cmd := &cobra.Command{
		Use:               "init [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a dir (default to the
current directory)`,
		RunE: func(_ *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			return run(log, api, in, dir, repositoryUrl, repositoryRef, recipeName)
		},
	}

	// Set flags
	cmd.Flags().StringVarP(&repositoryUrl, "repository", "o", "", "use repository")
	cmd.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	cmd.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")

	return cmd
}

func run(log *slog.Logger, api *api.Api, in ui.Input, dir, repositoryUrl, repositoryRef, recipeName string) error {
	// Get project finder
	projectFinder := api.NewProjectFinder()

	// Check already existing project
	log.Info("finding project…")
	if projectFinder.Find(dir) {
		return &app.AlreadyExistingProjectError{Dir: dir}
	}

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

	var recipe app.Recipe

	if recipeName != "" {
		// Load recipe by flag
		log.Info("loading recipe…")
		if recipe, err = recipeLoader.Load(repository, recipeName); err != nil {
			return err
		}
	} else {
		// Select recipe
		log.Info("loading recipes…")
		recipes, err := recipeLoader.LoadAll(repository)
		if err != nil {
			return err
		}

		form, err := NewUiRecipeListForm(recipes, &recipe)
		if err != nil {
			return err
		}

		if err := in.ListForm(
			"Please, select a recipe…",
			form,
		); err != nil {
			return err
		}
	}

	// Recipe vars
	vars := recipe.Vars()

	// Recipe options
	if len(recipe.Options()) > 0 {
		form, err := NewUiRecipeOptionsForm(recipe, &vars)
		if err != nil {
			return err
		}

		if err := in.Form(
			fmt.Sprintf("Please, enter \"%s\" recipe options…", recipe.Name()),
			form,
		); err != nil {
			return err
		}
	}

	// Create project
	log.Info("creating project…")
	project, err := api.NewProjectCreator().Create(dir, recipe, vars)
	if err != nil {
		return err
	}

	// Sync project
	log.Info("syncing project…")
	return api.NewProjectSyncer().Sync(project)
}
