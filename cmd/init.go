package cmd

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"gitlab.com/tslocum/cview"
	"io/fs"
	"manala/binder"
	"manala/config"
	"manala/loaders"
	"manala/logger"
	"manala/models"
	"manala/syncer"
	"manala/validator"
	"os"
	"path/filepath"
)

type InitCmd struct {
	Log              *logger.Logger
	Conf             *config.Config
	RepositoryLoader loaders.RepositoryLoaderInterface
	RecipeLoader     loaders.RecipeLoaderInterface
	ProjectLoader    loaders.ProjectLoaderInterface
	TemplateManager  models.TemplateManagerInterface
	Sync             *syncer.Syncer
	Assets           fs.ReadFileFS
}

func (cmd *InitCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:     "init [dir]",
		Aliases: []string{"in"},
		Short:   "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			// Get directory from first command arg
			dir := "."
			if len(args) != 0 {
				dir = args[0]
			}

			flags := command.Flags()

			cmd.Conf.BindRepositoryFlag(flags.Lookup("repository"))

			recName, _ := flags.GetString("recipe")

			return cmd.Run(dir, recName)
		},
	}

	flags := command.Flags()

	flags.StringP("repository", "o", "", "use repository source")
	flags.StringP("recipe", "i", "", "use recipe name")

	return command
}

func (cmd *InitCmd) Run(dir string, recName string) error {
	// Ensure directory exists
	if dir != "." {
		stat, err := os.Stat(dir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				cmd.Log.DebugWithField("Creating project directory...", "dir", dir)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("error creating project directory: %v", err)
				}
				cmd.Log.InfoWithField("Project directory created", "dir", dir)
			} else {
				return fmt.Errorf("error getting project directory stat: %v", err)
			}
		} else if !stat.IsDir() {
			return fmt.Errorf("project directory invalid: %s", dir)
		}
	}

	// Ensure no project already exists
	prjManifest, _ := cmd.ProjectLoader.Find(dir, false)
	if prjManifest != nil {
		return fmt.Errorf("project already exists: %s", dir)
	}

	// Load repository
	repo, err := cmd.RepositoryLoader.Load(cmd.Conf.Repository())
	if err != nil {
		return err
	}

	// Load recipe...
	var rec models.RecipeInterface
	if recName != "" {
		// ...from name if given
		rec, err = cmd.RecipeLoader.Load(recName, repo)
		if err != nil {
			return err
		}
	} else {
		// ...from recipe list
		rec, err = cmd.runRecipeListApplication(cmd.RecipeLoader, repo)
		if err != nil {
			return err
		}
	}

	// Vars
	vars := rec.Vars()

	// Use recipe options form if any
	if len(rec.Options()) > 0 {
		if err := cmd.runRecipeOptionsFormApplication(rec, vars); err != nil {
			return err
		}
	}

	// Template
	template, err := cmd.TemplateManager.NewRecipeTemplate(rec)
	if err != nil {
		return err
	}

	if rec.Template() != "" {
		// Load template from recipe
		if err := template.ParseFile(rec.Template()); err != nil {
			return err
		}
	} else {
		// Load default template from embedded assets
		text, _ := cmd.Assets.ReadFile("assets/" + models.ProjectManifestFile + ".tmpl")
		if err := template.Parse(string(text)); err != nil {
			return err
		}
	}

	// Create project manifest
	prjManifest, err = os.Create(filepath.Join(dir, models.ProjectManifestFile))
	if err != nil {
		return err
	}
	defer prjManifest.Close()

	if err := template.Execute(prjManifest, vars); err != nil {
		return err
	}

	prj, err := cmd.ProjectLoader.Load(prjManifest, "", "")
	if err != nil {
		return err
	}

	// Validate project
	if err := validator.ValidateProject(prj); err != nil {
		return err
	}

	cmd.Log.Info("Project validated")

	// Sync project
	if err := cmd.Sync.SyncProject(prj); err != nil {
		return err
	}

	cmd.Log.Info("Project synced")

	return nil
}

func (cmd *InitCmd) runRecipeListApplication(recLoader loaders.RecipeLoaderInterface, repo models.RepositoryInterface) (models.RecipeInterface, error) {
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

func (cmd *InitCmd) runRecipeOptionsFormApplication(rec models.RecipeInterface, vars map[string]interface{}) error {
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
	frame.AddText("Please, enter \""+rec.Name()+"\" recipe options...", true, cview.AlignLeft, tcell.ColorAqua)

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
	bndr, err2 := binder.NewRecipeFormBinder(rec)
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
			_ = bndr.Apply(vars)
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
