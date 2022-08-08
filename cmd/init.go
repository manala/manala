package cmd

import (
	"code.rocketnine.space/tslocum/cview"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"manala/app"
	"manala/internal"
	internalBinder "manala/internal/binder"
	internalConfig "manala/internal/config"
	internalErrors "manala/internal/errors"
	internalLog "manala/internal/log"
	internalValidator "manala/internal/validator"
	"path/filepath"
)

func newInitCmd(config *internalConfig.Config, logger *internalLog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "init [path]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a path (default to the current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// App
			manala := app.New(config, logger)

			// Get flags
			repositoryPath, _ := cmd.Flags().GetString("repository")
			recipeName, _ := cmd.Flags().GetString("recipe")

			// Get args
			path := filepath.Clean(append(args, "")[0])

			// Ensure no already existing project
			if manifest, err := manala.ProjectManifest(path); true {
				if manifest != nil {
					return internalErrors.New("already existing project").WithField("path", path)
				}
				var _err *internal.NotFoundProjectManifestError
				if !errors.As(err, &_err) {
					return err
				}
			}

			// Load repository
			repository, err := manala.Repository(repositoryPath)
			if err != nil {
				return err
			}

			// Load recipe...
			var recipe *internal.Recipe
			if recipeName != "" {
				// ...from name if passed as argument
				recipe, err = repository.LoadRecipe(recipeName)
				if err != nil {
					return err
				}
			} else {
				// ...from recipe list tui
				recipe, err = initRecipeListApplication(repository)
				if err != nil {
					return err
				}
			}

			// Create project
			project := recipe.NewProject(path)

			// Recipe options tui form (if any)
			if len(recipe.Options()) > 0 {
				if err := initRecipeOptionsFormApplication(recipe, project.Manifest()); err != nil {
					return err
				}
			}

			// Write manifest content
			if err := project.ManifestTemplate().Write(project.Manifest()); err != nil {
				return err
			}

			// Load manifest
			if err := project.Manifest().Load(); err != nil {
				return err
			}

			// Save manifest
			if err := project.Manifest().Save(); err != nil {
				return err
			}

			// Sync project
			return manala.SyncProject(project)
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")

	return cmd
}

func initRecipeListApplication(recipeWalker internal.RecipeWalkerInterface) (*internal.Recipe, error) {
	// Application
	application := cview.NewApplication()
	application.EnableMouse(true)

	var err error

	// List
	list := cview.NewList()
	list.SetPadding(0, 0, 1, 0)
	list.SetScrollBarVisibility(cview.ScrollBarAlways)
	list.SetDoneFunc(func() {
		err = fmt.Errorf("operation cancelled")
		application.Stop()
	})

	var recipe *internal.Recipe

	// Walk into recipes
	if err2 := recipeWalker.WalkRecipes(func(_recipe *internal.Recipe) {
		item := cview.NewListItem(" " + _recipe.Name() + " ")
		item.SetSecondaryText("   " + _recipe.Description())
		item.SetSelectedFunc(func() {
			recipe = _recipe
			application.Stop()
		})
		list.AddItem(item)
	}); err2 != nil {
		return nil, err2
	}

	frame := cview.NewFrame(list)
	frame.SetBorders(1, 1, 1, 1, 1, 1)
	frame.AddText("Please, select a recipe...", true, cview.AlignLeft, tcell.ColorAqua)

	application.SetRoot(frame, true)
	application.SetFocus(frame)
	if err2 := application.Run(); err2 != nil {
		return nil, err2
	}

	if err != nil {
		return nil, err
	}

	if recipe == nil {
		return nil, fmt.Errorf("unknown error")
	}

	return recipe, nil
}

func initRecipeOptionsFormApplication(recipe *internal.Recipe, manifest *internal.ProjectManifest) error {
	// Application
	application := cview.NewApplication()
	application.EnableMouse(true)

	panels := cview.NewPanels()

	var err error

	// Form panel
	form := cview.NewForm()
	form.SetPadding(0, 0, 1, 0)
	form.SetItemPadding(0)
	form.SetCancelFunc(func() {
		err = fmt.Errorf("operation cancelled")
		application.Stop()
	})

	frame := cview.NewFrame(form)
	frame.SetBorders(1, 1, 1, 1, 1, 1)
	frame.AddText("Please, enter \""+recipe.Name()+"\" recipe options...", true, cview.AlignLeft, tcell.ColorAqua)

	panels.AddPanel("form", frame, true, true)

	// Modal panel
	modal := cview.NewModal()
	modal.SetBorderColor(tcell.ColorRed)
	modal.SetTextColor(tcell.ColorWhite)
	modal.AddButtons([]string{"Ok"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		panels.HidePanel("modal")
	})

	panels.AddPanel("modal", modal, false, false)

	// Recipe form binder
	binder, _err := internalBinder.NewRecipeFormBinder(recipe.Options())
	if _err != nil {
		return _err
	}

	binder.BindForm(form)

	form.AddButton("Apply", func() {
		// Validate
		valid := true
		for _, bind := range binder.Binds() {
			if _err, _errs, ok := internalValidator.Validate(bind.Option.Schema, bind.Value); !ok {
				if _err != nil {
					err = _err
					application.Stop()
				} else {
					valid = false
					modal.SetText(bind.Option.Label + _errs[0].Error())
					panels.ShowPanel("modal")
					form.SetFocus(bind.ItemIndex)
				}
				break
			}
		}
		if valid && (err == nil) {
			// Apply values
			_ = binder.Apply(manifest)
			application.Stop()
		}
	})

	application.SetRoot(panels, true)
	application.SetFocus(panels)
	if err3 := application.Run(); err3 != nil {
		return err3
	}

	if err != nil {
		return err
	}

	return nil
}
