package cmd

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/apex/log"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/app"
	"manala/binder"
	"manala/loaders"
	"manala/models"
	"manala/validator"
)

type InitCmd struct{}

func (cmd *InitCmd) Command(config *viper.Viper, logger *log.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:     "init [dir]",
		Aliases: []string{"in"},
		Short:   "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			// Config
			_ = config.BindPFlags(command.PersistentFlags())

			// App
			manala := app.New(
				app.WithConfig(config),
				app.WithLogger(logger),
			)

			// Get directory from first command arg
			dir := "."
			if len(args) != 0 {
				dir = args[0]
			}

			// Flags
			flags := command.Flags()
			recName, _ := flags.GetString("recipe")

			// Command
			return manala.Init(
				cmd.runRecipeListApplication,
				cmd.runRecipeOptionsFormApplication,
				dir,
				recName,
			)
		},
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("repository", "o", "", "use repository source")

	// Flags
	flags := command.Flags()
	flags.StringP("recipe", "i", "", "use recipe name")

	return command
}

func (cmd *InitCmd) runRecipeListApplication(recipeLoader loaders.RecipeLoaderInterface, repo models.RepositoryInterface) (models.RecipeInterface, error) {
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

	var recipe models.RecipeInterface

	// Walk into recipes
	if err2 := recipeLoader.Walk(repo, func(rec models.RecipeInterface) {
		listItem := cview.NewListItem(" " + rec.Name() + " ")
		listItem.SetSecondaryText("   " + rec.Description())
		listItem.SetSelectedFunc(func() {
			recipe = rec
			application.Stop()
		})
		list.AddItem(listItem)
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

func (cmd *InitCmd) runRecipeOptionsFormApplication(rec models.RecipeInterface, vars map[string]interface{}) error {
	// Application
	application := cview.NewApplication()
	application.EnableMouse(true)

	appPanels := cview.NewPanels()

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
					application.Stop()
				}
				break
			}
		}
		if valid && (err == nil) {
			// Apply values
			_ = bndr.Apply(vars)
			application.Stop()
		}
	})

	application.SetRoot(appPanels, true)
	application.SetFocus(appPanels)
	if err3 := application.Run(); err3 != nil {
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
