package init

import (
	"context"
	"fmt"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/internal/ui"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewCmd(log *slog.Logger, api *api.API, input ui.Input) *cobra.Command {
	// Flags
	var repositoryURL, repositoryRef, recipeName string

	// Command
	cmd := &cobra.Command{
		Use:               "init [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a dir (default to the
current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			// Context
			ctx := cmd.Context()
			ctx = app.WithRepositoryURL(ctx, repositoryURL)
			ctx = app.WithRepositoryRef(ctx, repositoryRef)
			ctx = app.WithRecipeName(ctx, recipeName)

			return run(ctx, log, api, input, dir)
		},
	}

	// Set flags
	cmd.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	cmd.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	cmd.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")

	return cmd
}

func run(ctx context.Context, log *slog.Logger, api *api.API, input ui.Input, dir string) error {
	// Get project finder
	projectFinder := api.NewProjectFinder()

	// Check already existing project
	log.Info("finding project…")

	if projectFinder.Find(dir) {
		return &app.AlreadyExistingProjectError{Dir: dir}
	}

	// Get repository loader
	repositoryLoader := api.NewRepositoryLoader(ctx)

	// Load repository
	log.Info("loading repository…")

	repository, err := repositoryLoader.Load("")
	if err != nil {
		return err
	}

	// Get recipe loader
	recipeLoader := api.NewRecipeLoader(ctx)

	var recipe app.Recipe

	if _, ok := app.RecipeName(ctx); ok {
		// Load recipe by context
		log.Info("loading recipe…")

		if recipe, err = recipeLoader.Load(repository, ""); err != nil {
			return err
		}
	} else {
		// Select recipe
		log.Info("loading recipes…")

		recipes, err := recipeLoader.LoadAll(repository)
		if err != nil {
			return err
		}

		form, err := NewUIRecipeListForm(recipes, &recipe)
		if err != nil {
			return err
		}

		if err := input.ListForm(
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
		form, err := NewUIRecipeOptionsForm(recipe, &vars)
		if err != nil {
			return err
		}

		if err := input.Form(
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
