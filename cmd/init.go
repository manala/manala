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
	app.EnableMouse()

	var error error

	// List
	list := cview.NewList()
	list.
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetTitleAlign(cview.AlignLeft).
		SetTitle(" Select recipe ")
	list.
		SetScrollBarVisibility(cview.ScrollBarAlways).
		SetDoneFunc(func() {
			error = fmt.Errorf("operation cancelled")
			app.Stop()
		})

	var recipe models.RecipeInterface

	// Walk into recipes
	if err := recLoader.Walk(repo, func(rec models.RecipeInterface) {
		list.AddItem(rec.Name(), rec.Description(), 0, func() {
			recipe = rec
			app.Stop()
		})
	}); err != nil {
		return nil, err
	}

	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
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
	app.EnableMouse()

	appPages := cview.NewPages()

	var error error

	// Form page
	form := cview.NewForm()
	form.
		SetBorder(true).
		SetTitleAlign(cview.AlignLeft).
		SetTitle(" Enter \"" + prj.Recipe().Name() + "\" options ")
	form.
		SetItemPadding(0).
		SetCancelFunc(func() {
			error = fmt.Errorf("operation cancelled")
			app.Stop()
		})

	appPages.AddPage("form", form, true, true)

	// Modal page
	modal := cview.NewModal()
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
		//PrimitiveBackgroundColor:    tcell.ColorBlack,
		PrimitiveBackgroundColor:    tcell.ColorDefault,
		ContrastBackgroundColor:     tcell.ColorBlue,
		MoreContrastBackgroundColor: tcell.ColorGreen,
		BorderColor:                 tcell.ColorWhite,
		TitleColor:                  tcell.ColorWhite,
		GraphicsColor:               tcell.ColorWhite,
		//PrimaryTextColor:            tcell.ColorWhite,
		PrimaryTextColor:            tcell.ColorDefault,
		SecondaryTextColor:          tcell.ColorYellow,
		TertiaryTextColor:           tcell.ColorGreen,
		InverseTextColor:            tcell.ColorBlue,
		ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
		ScrollBarColor:              tcell.ColorWhite,
	}
}