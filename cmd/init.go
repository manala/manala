package cmd

import (
	"code.rocketnine.space/tslocum/cview"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"manala/core"
	"manala/core/application"
	internalBinder "manala/internal/binder"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
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
			// Application
			app := application.NewApplication(config, logger)

			// Get flags
			repoPath, _ := cmd.Flags().GetString("repository")
			recName, _ := cmd.Flags().GetString("recipe")

			// Get args
			path := filepath.Clean(append(args, "")[0])

			// Ensure no already existing project
			if manifest, err := app.ProjectManifest(path); true {
				if manifest != nil {
					return internalReport.NewError(fmt.Errorf("already existing project")).
						WithField("path", path)
				}
				var _notFoundProjectManifestError *core.NotFoundProjectManifestError
				if !errors.As(err, &_notFoundProjectManifestError) {
					return err
				}
			}

			// Load repository
			repo, err := app.Repository(repoPath)
			if err != nil {
				return err
			}

			// Create project
			proj, err := app.CreateProject(
				path,
				repo,
				// Recipe selector
				func(recWalker core.RecipeWalker) (core.Recipe, error) {
					// From argument
					if recName != "" {
						var rec core.Recipe
						if err := recWalker.WalkRecipes(func(_rec core.Recipe) {
							if _rec.Name() == recName {
								rec = _rec
							}
						}); err != nil {
							return nil, err
						}
						if rec == nil {
							return nil, fmt.Errorf("recipe not found")
						}
						return rec, nil
					}

					// From tui list
					return initRecipeListApplication(recWalker)
				},
				// Options selector
				func(rec core.Recipe, options []core.RecipeOption) error {
					if len(options) > 0 {
						// From tui form
						return initRecipeOptionsFormApplication(rec, options)
					}

					return nil
				},
			)
			if err != nil {
				return err
			}

			// Sync project
			return app.SyncProject(proj)
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")

	return cmd
}

func initRecipeListApplication(recWalker core.RecipeWalker) (core.Recipe, error) {
	// Application
	app := cview.NewApplication()
	app.EnableMouse(true)

	var err error

	// List
	list := cview.NewList()
	list.SetPadding(0, 0, 1, 0)
	list.SetScrollBarVisibility(cview.ScrollBarAlways)
	list.SetDoneFunc(func() {
		err = fmt.Errorf("operation cancelled")
		app.Stop()
	})

	var rec core.Recipe

	// Walk into recipes
	if err2 := recWalker.WalkRecipes(func(_rec core.Recipe) {
		item := cview.NewListItem(" " + _rec.Name() + " ")
		item.SetSecondaryText("   " + _rec.Description())
		item.SetSelectedFunc(func() {
			rec = _rec
			app.Stop()
		})
		list.AddItem(item)
	}); err2 != nil {
		return nil, err2
	}

	frame := cview.NewFrame(list)
	frame.SetBorders(1, 1, 1, 1, 1, 1)
	frame.AddText("Please, select a recipe...", true, cview.AlignLeft, tcell.ColorAqua)

	app.SetRoot(frame, true)
	app.SetFocus(frame)
	if err2 := app.Run(); err2 != nil {
		return nil, err2
	}

	if err != nil {
		return nil, err
	}

	if rec == nil {
		return nil, fmt.Errorf("unknown error")
	}

	return rec, nil
}

func initRecipeOptionsFormApplication(rec core.Recipe, options []core.RecipeOption) error {
	// Application
	app := cview.NewApplication()
	app.EnableMouse(true)

	panels := cview.NewPanels()

	var err error

	// Form panel
	form := cview.NewForm()
	form.SetPadding(0, 0, 1, 0)
	form.SetItemPadding(0)
	form.SetCancelFunc(func() {
		err = fmt.Errorf("operation cancelled")
		app.Stop()
	})

	frame := cview.NewFrame(form)
	frame.SetBorders(1, 1, 1, 1, 1, 1)
	frame.AddText("Please, enter \""+rec.Name()+"\" recipe options...", true, cview.AlignLeft, tcell.ColorAqua)

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
	binder, _err := internalBinder.NewRecipeFormBinder(options)
	if _err != nil {
		return _err
	}

	binder.BindForm(form)

	form.AddButton("Apply", func() {
		// Validate
		valid := true
		for _, bind := range binder.Binds() {
			validation, _err := gojsonschema.Validate(
				gojsonschema.NewGoLoader(bind.Option.Schema()),
				gojsonschema.NewGoLoader(bind.Value),
			)

			if _err != nil {
				err = _err
				app.Stop()
				break
			}

			if !validation.Valid() {
				valid = false
				modal.SetText(bind.Option.Label() + validation.Errors()[0].String())
				panels.ShowPanel("modal")
				form.SetFocus(bind.ItemIndex)
				break
			}
		}
		if valid && (err == nil) {
			// Apply values
			_ = binder.Apply()
			app.Stop()
		}
	})

	app.SetRoot(panels, true)
	app.SetFocus(panels)
	if err3 := app.Run(); err3 != nil {
		return err3
	}

	if err != nil {
		return err
	}

	return nil
}
