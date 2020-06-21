package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gdamore/tcell"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/tslocum/cview"
	"manala/binder"
	"manala/loaders"
	"manala/models"
	"manala/syncer"
	"manala/validator"
)

// InitCmd represents the init command
func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"in"},
		Short:   "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in the current directory`,
		Run:  initRun,
		Args: cobra.NoArgs,
	}

	return cmd
}

func initRun(cmd *cobra.Command, args []string) {
	// Loaders
	repoLoader := loaders.NewRepositoryLoader(viper.GetString("cache_dir"))
	recLoader := loaders.NewRecipeLoader()
	prjLoader := loaders.NewProjectLoader(repoLoader, recLoader, viper.GetString("repository"))

	// Ensure project is not yet initialized by checking configuration file existence
	cfgFile, _ := prjLoader.ConfigFile(viper.GetString("dir"))
	if cfgFile != nil {
		log.Fatal("Project already initialized")
	}

	// Load repository
	repo, err := repoLoader.Load(viper.GetString("repository"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Recipe list application
	rec, err := initRecipeListApplication(recLoader, repo)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Project
	prj := models.NewProject(
		viper.GetString("dir"),
		rec,
	)

	if rec.HasOptions() {
		// Project form application
		if err := initProjectFormApplication(prj); err != nil {
			log.Fatal(err.Error())
		}
	}

	// Sync project
	if err := syncer.SyncProject(prj); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Project synced")
}

func initRecipeListApplication(recLoader loaders.RecipeLoaderInterface, repo models.RepositoryInterface) (models.RecipeInterface, error) {
	// Application
	app := cview.NewApplication()
	app.EnableMouse(true)

	var error error

	// List
	list := cview.NewList()
	list.SetBorderPadding(0, 0, 1, 0)
	list.
		SetScrollBarVisibility(cview.ScrollBarAlways).
		SetDoneFunc(func() {
			error = fmt.Errorf("operation cancelled")
			app.Stop()
		})

	var recipe models.RecipeInterface

	// Walk into recipes
	if err := recLoader.Walk(repo, func(rec models.RecipeInterface) {
		list.AddItem(" " + rec.Name() + " ", "   " + rec.Description(), 0, func() {
			recipe = rec
			app.Stop()
		})
	}); err != nil {
		return nil, err
	}

	frame := cview.NewFrame(list).
		SetBorders(1, 1, 1, 1, 1, 1).
		AddText("Please, select a recipe...", true, cview.AlignLeft, tcell.ColorAqua)

	if err := app.SetRoot(frame, true).SetFocus(frame).Run(); err != nil {
		return nil, err
	}

	if error != nil {
		return nil, error
	}

	if recipe == nil {
		return nil, fmt.Errorf("unknown error")
	}

	return recipe, nil
}

func initProjectFormApplication(prj models.ProjectInterface) error {
	// Application
	app := cview.NewApplication()
	app.EnableMouse(true)

	appPages := cview.NewPages()

	var error error

	// Form page
	form := cview.NewForm()
	form.SetBorderPadding(0, 0, 1, 0)
	form.
		SetItemPadding(0).
		SetCancelFunc(func() {
			error = fmt.Errorf("operation cancelled")
			app.Stop()
		})

	frame := cview.NewFrame(form).
		SetBorders(1, 1, 1, 1, 1, 1).
		AddText("Please, enter \"" + prj.Recipe().Name() + "\" recipe options...", true, cview.AlignLeft, tcell.ColorAqua)

	appPages.AddPage("form", frame, true, true)

	// Modal page
	modal := cview.NewModal()
	modal.SetBorderColor(tcell.ColorRed)
	modal.SetTextColor(tcell.ColorWhite)
	modal.
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			appPages.HidePage("modal")
		})

	appPages.AddPage("modal", modal, false, false)

	// Recipe form binder
	bndr, err := binder.NewRecipeFormBinder(prj.Recipe())
	if err != nil {
		return err
	}

	bndr.BindForm(form)

	form.AddButton("Apply", func() {
		// Validate
		valid := true
		for _, bnd := range bndr.Binds() {
			err := validator.ValidateValue(bnd.Value, bnd.Option.Schema)
			if err != nil {
				if err, ok := err.(*validator.ValueValidationError); ok {
					valid = false
					modal.SetText(bnd.Option.Label + err.Error())
					appPages.ShowPage("modal")
					form.SetFocus(bnd.ItemIndex)
				} else {
					error = err
					app.Stop()
				}
				break
			}
		}
		if valid && err == nil {
			// Apply values
			bndr.ApplyValues(prj.Vars())
			app.Stop()
		}
	})

	if err := app.SetRoot(appPages, true).SetFocus(appPages).Run(); err != nil {
		return err
	}

	if error != nil {
		return error
	}

	return nil
}

func init() {
	cview.Styles = cview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorBlack,
		ContrastBackgroundColor:     tcell.ColorBlack,
		MoreContrastBackgroundColor: tcell.ColorWhite,
		BorderColor:                 tcell.ColorRed,
		TitleColor:                  tcell.ColorAqua,
		GraphicsColor:               tcell.ColorWhite,
		PrimaryTextColor:            tcell.ColorLime,
		SecondaryTextColor:          tcell.ColorWhite,
		TertiaryTextColor:           tcell.ColorWhite,
		InverseTextColor:            tcell.ColorYellow,
		ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
		ScrollBarColor:              tcell.ColorWhite,
	}
}