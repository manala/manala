package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/app/config"
	"manala/internal/ui"
	"path/filepath"
)

func newInitCmd(config config.Config, log *slog.Logger, out ui.Output, in ui.Input) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "init [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a dir (default to the current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Api options
			var apiOptions []api.Option

			// Flag - Repository url
			if cmd.Flags().Changed("repository") {
				repositoryUrl, _ := cmd.Flags().GetString("repository")
				apiOptions = append(apiOptions, api.WithRepositoryUrl(repositoryUrl))
			}

			// Flag - Repository ref
			if cmd.Flags().Changed("ref") {
				repositoryRef, _ := cmd.Flags().GetString("ref")
				apiOptions = append(apiOptions, api.WithRepositoryRef(repositoryRef))
			}

			// Flag - Recipe name
			if cmd.Flags().Changed("recipe") {
				recipeName, _ := cmd.Flags().GetString("recipe")
				apiOptions = append(apiOptions, api.WithRecipeName(recipeName))
			}

			// Api
			api := api.New(
				config,
				log,
				out,
				apiOptions...,
			)

			// Get args
			dir := filepath.Clean(append(args, "")[0])

			// Check already existing project
			log.Debug("check already existing project…", "dir", dir)
			if api.IsProject(dir) {
				return &app.AlreadyExistingProjectError{Dir: dir}
			}

			// Load preceding repository
			log.Debug("load preceding repository…")
			repository, err := api.LoadPrecedingRepository()
			if err != nil {
				return err
			}

			var recipe app.Recipe

			// Try loading preceding recipe
			recipe, err = api.LoadPrecedingRecipe(repository)

			if err != nil {
				var _unprocessableRecipeNameError *app.UnprocessableRecipeNameError
				if !errors.As(err, &_unprocessableRecipeNameError) {
					return err
				}

				log.Debug("unable to load preceding recipe")

				recipes, err := api.RepositoryRecipes(repository)
				if err != nil {
					return err
				}

				// Select recipe
				log.Debug("select recipe…")

				form, err := api.NewUiRecipeListForm(recipes, &recipe)
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

				// Set recipe options
				log.Debug("set recipe options…")

				form, err := api.NewUiRecipeOptionsForm(recipe, &vars)
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
			project, err := api.CreateProject(dir, recipe, vars)
			if err != nil {
				return err
			}

			// Sync project
			return api.SyncProject(project)
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().String("ref", "", "use repository ref")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")

	return cmd
}
