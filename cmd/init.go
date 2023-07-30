package cmd

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"log/slog"
	"manala/app/interfaces"
	"manala/core/application"
	"manala/internal/binder"
	"manala/internal/ui/output"
	"path/filepath"
)

func newInitCmd(conf interfaces.Config, log *slog.Logger, out output.Output) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "init [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a dir (default to the current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Application options
			var appOptions []application.Option

			// Flag - Repository url
			if cmd.Flags().Changed("repository") {
				repoUrl, _ := cmd.Flags().GetString("repository")
				appOptions = append(appOptions, application.WithRepositoryUrl(repoUrl))
			}

			// Flag - Repository ref
			if cmd.Flags().Changed("ref") {
				repoRef, _ := cmd.Flags().GetString("ref")
				appOptions = append(appOptions, application.WithRepositoryRef(repoRef))
			}

			// Flag - Recipe name
			if cmd.Flags().Changed("recipe") {
				recName, _ := cmd.Flags().GetString("recipe")
				appOptions = append(appOptions, application.WithRecipeName(recName))
			}

			// Application
			app := application.NewApplication(
				conf,
				log,
				out,
				appOptions...,
			)

			// Get args
			dir := filepath.Clean(append(args, "")[0])

			// Create project
			proj, err := app.CreateProject(
				dir,
				// Recipe selector
				func(recWalker func(walker func(rec interfaces.Recipe) error) error) (interfaces.Recipe, error) {
					// From tui list
					return initRecipeListApplication(recWalker)
				},
				// Options selector
				func(rec interfaces.Recipe, options []interfaces.RecipeOption) error {
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
	cmd.Flags().String("ref", "", "use repository ref")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")

	return cmd
}

func initRecipeListApplication(recWalker func(walker func(rec interfaces.Recipe) error) error) (interfaces.Recipe, error) {
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

	var rec interfaces.Recipe

	// Walk into recipes
	if err2 := recWalker(func(_rec interfaces.Recipe) error {
		item := cview.NewListItem(" " + _rec.Name() + " ")
		item.SetSecondaryText("   " + _rec.Description())
		item.SetSelectedFunc(func() {
			rec = _rec
			app.Stop()
		})
		list.AddItem(item)
		return nil
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

func initRecipeOptionsFormApplication(rec interfaces.Recipe, options []interfaces.RecipeOption) error {
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
	binder, _err := binder.NewRecipeFormBinder(options)
	if _err != nil {
		return _err
	}

	binder.BindForm(form)

	form.AddButton("Apply", func() {
		// Validate
		valid := true
		for _, bind := range binder.Binds() {
			val, _err := gojsonschema.Validate(
				gojsonschema.NewGoLoader(bind.Option.Schema()),
				gojsonschema.NewGoLoader(bind.Value),
			)

			if _err != nil {
				err = _err
				app.Stop()
				break
			}

			if !val.Valid() {
				valid = false
				modal.SetText(bind.Option.Label() + val.Errors()[0].String())
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
