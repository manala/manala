package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/tslocum/cview"
	"manala/binder"
	"manala/loaders"
	"manala/models"
	"manala/syncer"
	"manala/validator"
	"os"
)

// InitCmd represents the init command
func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [dir]",
		Aliases: []string{"in"},
		Short:   "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE:              initRun,
	}

	addRepositoryFlag(cmd, "use repository")
	addRecipeFlag(cmd, "use recipe")

	return cmd
}

func initRun(cmd *cobra.Command, args []string) error {
	// Loaders
	repoLoader := loaders.NewRepositoryLoader(
		viper.GetString("cache_dir"),
		viper.GetString("repository"),
	)
	recLoader := loaders.NewRecipeLoader()
	prjLoader := loaders.NewProjectLoader(repoLoader, recLoader, "", "")

	// Directory
	dir := "."
	if len(args) != 0 {
		// Get directory from first command arg
		dir = args[0]
		// Ensure directory exists
		stat, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				log.WithField("dir", dir).Debug("Creating project directory...")
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("error creating project directory: %v", err)
				}
				log.WithField("dir", dir).Info("Project directory created")
			} else {
				return fmt.Errorf("error getting project directory stat: %v", err)
			}
		} else if !stat.IsDir() {
			return fmt.Errorf("project directory invalid: %s", dir)
		}
	}

	// Ensure no project already exists
	prjFile, _ := prjLoader.Find(dir, false)
	if prjFile != nil {
		return fmt.Errorf("project already exists: %s", dir)
	}

	// Load repository
	repoSrc, _ := cmd.Flags().GetString("repository")
	repo, err := repoLoader.Load(repoSrc)
	if err != nil {
		return err
	}

	// Recipe
	var rec models.RecipeInterface

	// From command flag
	recName, _ := cmd.Flags().GetString("recipe")
	if recName != "" {
		rec, err = recLoader.Load(recName, repo)
		if err != nil {
			return err
		}
	} else {
		// From recipe list application
		rec, err = initRecipeListApplication(recLoader, repo)
		if err != nil {
			return err
		}
	}

	// Project
	prj := models.NewProject(dir, rec)

	if rec.HasOptions() {
		// Project form application
		if err := initProjectFormApplication(prj); err != nil {
			return err
		}
	}

	// Sync project
	if err := syncer.SyncProject(prj); err != nil {
		return err
	}

	log.Info("Project synced")

	return nil
}

func initRecipeListApplication(recLoader loaders.RecipeLoaderInterface, repo models.RepositoryInterface) (models.RecipeInterface, error) {
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

	var recipe models.RecipeInterface

	// Walk into recipes
	if err2 := recLoader.Walk(repo, func(rec models.RecipeInterface) {
		listItem := cview.NewListItem(" " + rec.Name() + " ")
		listItem.SetSecondaryText("   " + rec.Description())
		listItem.SetSelectedFunc(func() {
			recipe = rec
			app.Stop()
		})
		list.AddItem(listItem)
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

	if recipe == nil {
		return nil, fmt.Errorf("unknown error")
	}

	return recipe, nil
}

func initProjectFormApplication(prj models.ProjectInterface) error {
	// Application
	app := cview.NewApplication()
	app.EnableMouse(true)

	appPanels := cview.NewPanels()

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
	frame.AddText("Please, enter \""+prj.Recipe().Name()+"\" recipe options...", true, cview.AlignLeft, tcell.ColorAqua)

	appPanels.AddPanel("form", frame, true, true)

	// Modal panel
	modal := cview.NewModal()
	modal.SetBorderColor(tcell.ColorRed)
	modal.SetTextColor(tcell.ColorWhite)
	modal.AddButtons([]string{"Ok"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		appPanels.HidePanel("modal")
	})

	appPanels.AddPanel("modal", modal, false, false)

	// Recipe form binder
	bndr, err2 := binder.NewRecipeFormBinder(prj.Recipe())
	if err2 != nil {
		return err2
	}

	bndr.BindForm(form)

	form.AddButton("Apply", func() {
		// Validate
		valid := true
		for _, bnd := range bndr.Binds() {
			err2 := validator.ValidateValue(bnd.Value, bnd.Option.Schema)
			if err2 != nil {
				if err3, ok := err2.(*validator.ValueValidationError); ok {
					valid = false
					modal.SetText(bnd.Option.Label + err3.Error())
					appPanels.ShowPanel("modal")
					form.SetFocus(bnd.ItemIndex)
				} else {
					err = err2
					app.Stop()
				}
				break
			}
		}
		if valid && (err == nil) {
			// Apply values
			_ = bndr.ApplyValues(prj.Vars())
			app.Stop()
		}
	})

	app.SetRoot(appPanels, true)
	app.SetFocus(appPanels)
	if err3 := app.Run(); err3 != nil {
		return err3
	}

	if err != nil {
		return err
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
